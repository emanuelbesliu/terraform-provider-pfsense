package pfsense

// BULK APPLY MODE (prototype v1, PFSenseBulkApply=1):
// Coalesces firewall rule creates and deletes into ONE php exec on the
// firewall: load config, N mutations, write_config(), filter_configure().
// Flush happens at the first Read after pending writes (read-after-write
// consistency within the same terraform apply) and via a size threshold so
// pure-destroy runs cannot lose their queue. Floating rules and unsupported
// shapes fall back to the classic per-item path.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

func BulkEnabled() bool { return os.Getenv("PFSenseBulkApply") == "1" }

type bulkOp struct {
	op      string // "create" | "update" | "delete"
	rule    map[string]string
	tracker string
	pending string // placeholder id returned to Terraform (create only)
}

type bulkState struct {
	mu         sync.Mutex
	ops        []bulkOp
	cache      map[string]string
	resolved   map[string]string // placeholder -> real tracker, from the flush that assigned it (1:1, order-based)
	rulesCache []FirewallRule
	fetchedAt  time.Time
}

var (
	bulkMu      sync.Mutex
	bulkStates  = map[*Client]*bulkState{}
	bulkCounter int64
)

func (pf *Client) bulk() *bulkState {
	bulkMu.Lock()
	defer bulkMu.Unlock()
	st, ok := bulkStates[pf]
	if !ok {
		st = &bulkState{cache: map[string]string{}, resolved: map[string]string{}}
		bulkStates[pf] = st
	}
	return st
}

func RuleSupportedForBulk(r FirewallRule) bool {
	return !r.Floating
}

func ruleToBulkCreate(r FirewallRule) map[string]string {
	m := map[string]string{
		"type":        r.Type,
		"interface":   r.Interface,
		"ipprotocol":  r.IPProtocol,
		"proto":       r.Protocol,
		"source":      r.SourceAddress,
		"srcport":     r.SourcePort,
		"destination": r.DestAddress,
		"dstport":     r.DestPort,
		"descr":       r.Description,
	}
	if r.SourceNot {
		m["source_not"] = "1"
	}
	if r.DestNot {
		m["dest_not"] = "1"
	}
	if r.Disabled {
		m["disabled"] = "1"
	}
	if r.Log {
		m["log"] = "1"
	}
	return m
}

// EnqueueCreate stores a create op; returns the placeholder tracker.
func (pf *Client) EnqueueCreate(r FirewallRule) string {
	st := pf.bulk()
	st.mu.Lock()
	defer st.mu.Unlock()
	ph := fmt.Sprintf("pending-%d", atomic.AddInt64(&bulkCounter, 1))
	st.ops = append(st.ops, bulkOp{op: "create", rule: ruleToBulkCreate(r), pending: ph})
	return ph
}

// EnqueueUpdate stores an in-place update op (tracker preserved on the box).
func (pf *Client) EnqueueUpdate(r FirewallRule) {
	st := pf.bulk()
	st.mu.Lock()
	defer st.mu.Unlock()
	st.ops = append(st.ops, bulkOp{op: "update", tracker: r.Tracker, rule: ruleToBulkCreate(r)})
}

// EnqueueDelete stores a delete-by-tracker op.
func (pf *Client) EnqueueDelete(tracker string) {
	st := pf.bulk()
	st.mu.Lock()
	defer st.mu.Unlock()
	st.ops = append(st.ops, bulkOp{op: "delete", tracker: tracker})
}

func (pf *Client) PendingCount() int {
	st := pf.bulk()
	st.mu.Lock()
	defer st.mu.Unlock()
	return len(st.ops)
}

// FlushPending applies all queued ops in one php exec and refreshes the
// description->tracker cache.
func (pf *Client) FlushPending(ctx context.Context) error {
	st := pf.bulk()
	st.mu.Lock()
	ops := st.ops
	st.ops = nil
	st.mu.Unlock()

	if len(ops) == 0 {
		return nil
	}

	type batch struct {
		Adds    []map[string]string `json:"adds"`
		Updates []map[string]string `json:"updates"`
		Deletes []string            `json:"deletes"`
	}
	b := batch{Adds: []map[string]string{}, Updates: []map[string]string{}, Deletes: []string{}}
	addPending := []string{} // placeholder per add, in the exact order the php loop assigns trackers
	for _, o := range ops {
		switch o.op {
		case "create":
			b.Adds = append(b.Adds, o.rule)
			addPending = append(addPending, o.pending)
		case "update":
			u := o.rule
			u["tracker"] = o.tracker
			b.Updates = append(b.Updates, u)
		case "delete":
			b.Deletes = append(b.Deletes, o.tracker)
		}
	}
	raw, _ := json.Marshal(b)
	b64 := base64.StdEncoding.EncodeToString(raw)

	script := `require_once("config.inc"); require_once("filter.inc"); global $config;
$b = json_decode(base64_decode('` + b64 + `'), true);
if (!isset($config['filter']['rule'])) { $config['filter']['rule'] = []; }
$mt = 0;
foreach ($config['filter']['rule'] as $rr0) { $t0v = (int)($rr0['tracker'] ?? 0); if ($t0v > $mt) { $mt = $t0v; } }
$ti = 0; $asg = array();
foreach ($b['adds'] as $r) {
	$ti++;
	$rule = array(
		'type' => $r['type'], 'interface' => $r['interface'],
		'ipprotocol' => $r['ipprotocol'],
		'source' => array('network' => $r['source']),
		'destination' => array('network' => $r['destination']),
		'descr' => $r['descr'],
	);
	if (isset($r['proto']) && "" !== (string)$r['proto']) { $rule['proto'] = $r['proto']; }
	if ("" !== (string)$r['srcport']) { $rule['source']['port'] = $r['srcport']; }
	if ("" !== (string)$r['dstport']) { $rule['destination']['port'] = $r['dstport']; }
	if (isset($r['log']) && in_array(strtolower((string)$r['log']), array('true', '1'), true)) { $rule['log'] = '1'; }
	if (isset($r['disabled']) && in_array(strtolower((string)$r['disabled']), array('true', '1'), true)) { $rule['disabled'] = '1'; }
	if (!empty($r['source_not'])) { $rule['source']['not'] = ''; }
	if (!empty($r['dest_not'])) { $rule['destination']['not'] = ''; }
	$rule['tracker'] = $mt + $ti;
	array_push($asg, strval($rule['tracker']));
	$config['filter']['rule'][] = $rule;
}
if (count($b['updates']) > 0) {
	$idx = array();
	foreach ($config['filter']['rule'] as $kk => $rr0) { $idx[(string)($rr0['tracker'] ?? '')] = $kk; }
	foreach ($b['updates'] as $u) {
		$k = $idx[$u['tracker']] ?? null;
		if ($k === null) { continue; }
		$rr = &$config['filter']['rule'][$k];
		$rr['type'] = $u['type']; $rr['interface'] = $u['interface'];
		$rr['ipprotocol'] = $u['ipprotocol'];
		if ("" !== (string)$u['proto']) { $rr['proto'] = $u['proto']; } else { unset($rr['proto']); }
		$rr['source'] = ("" === (string)$u['srcport']) ? array('network' => $u['source']) : array('network' => $u['source'], 'port' => $u['srcport']);
		if (!empty($u['source_not'])) { $rr['source']['not'] = ''; } else { unset($rr['source']['not']); }
		$rr['destination'] = ("" === (string)$u['dstport']) ? array('network' => $u['destination']) : array('network' => $u['destination'], 'port' => $u['dstport']);
		if (!empty($u['dest_not'])) { $rr['destination']['not'] = ''; } else { unset($rr['destination']['not']); }
		$rr['descr'] = $u['descr'];
		if (isset($u['log']) && in_array(strtolower((string)$u['log']), array('true', '1'), true)) { $rr['log'] = '1'; } else { unset($rr['log']); }
		if (isset($u['disabled']) && in_array(strtolower((string)$u['disabled']), array('true', '1'), true)) { $rr['disabled'] = '1'; } else { unset($rr['disabled']); }
		unset($rr);
	}
}
if (count($b['deletes']) > 0) {
	$dead = array_flip($b['deletes']);
	$out = array();
	foreach ($config['filter']['rule'] as $rr) {
		if (!isset($dead[(string)($rr['tracker'] ?? '')])) { $out[] = $rr; }
	}
	$config['filter']['rule'] = array_values($out);
}
write_config("pfsense-provider bulk apply");
filter_configure();
clear_subsystem_dirty('filter');
echo json_encode(array("asg" => array_values($asg)));`

	pf.mutexes.FirewallRule.Lock()
	var result struct {
		Asg []string `json:"asg"`
	}
	err := func() error {
		defer pf.mutexes.FirewallRule.Unlock()
		return pf.executePHPCommand(ctx, script, &result)
	}()
	if err != nil {
		return fmt.Errorf("bulk flush failed: %w", err)
	}

	// Deterministic placeholder resolution: the flush itself returned the tracker it
	// assigned to each add, in the same order - no description matching needed. This is
	// race-proof for duplicate descriptions, which the legacy desc-map cannot be.
	if len(result.Asg) == len(addPending) && len(addPending) > 0 {
		st.mu.Lock()
		if st.resolved == nil {
			st.resolved = map[string]string{}
		}
		for i, ph := range addPending {
			st.resolved[ph] = result.Asg[i]
		}
		st.mu.Unlock()
	}

	// refresh description -> tracker cache (only needed when creates added new rules)
	hasAddsOrUpdates := false
	for _, o := range ops {
		if o.op == "create" || o.op == "update" {
			hasAddsOrUpdates = true
			break
		}
	}
	if !hasAddsOrUpdates {
		return nil
	}
	rules, err := pf.GetFirewallRules(ctx)
	if err != nil {
		return err
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	st.cache = map[string]string{}
	st.rulesCache = *rules
	st.fetchedAt = time.Now()
	for _, rl := range *rules {
		if rl.Description != "" {
			st.cache[rl.Description] = rl.Tracker
		}
	}
	return nil
}

// EnsureTrackerCache fills description->tracker map from the live box if empty.
func (pf *Client) EnsureTrackerCache(ctx context.Context) error {
	st := pf.bulk()
	st.mu.Lock()
	fresh := len(st.rulesCache) > 0 && time.Since(st.fetchedAt) < 30*time.Second
	st.mu.Unlock()
	if fresh {
		return nil
	}
	rules, err := pf.GetFirewallRules(ctx)
	if err != nil {
		return err
	}
	st.mu.Lock()
	defer st.mu.Unlock()
	st.cache = map[string]string{}
	st.rulesCache = *rules
	st.fetchedAt = time.Now()
	for _, rl := range *rules {
		if rl.Description != "" {
			st.cache[rl.Description] = rl.Tracker
		}
	}
	return nil
}

// ResolvePlaceholder returns the real tracker a flush assigned to a "pending-" id.
func (pf *Client) ResolvePlaceholder(pending string) (string, bool) {
	st := pf.bulk()
	st.mu.Lock()
	defer st.mu.Unlock()
	t, ok := st.resolved[pending]
	return t, ok && t != ""
}

// LookupTracker resolves a description to its real tracker from the last flush.
func (pf *Client) LookupTracker(descr string) string {
	st := pf.bulk()
	st.mu.Lock()
	defer st.mu.Unlock()
	return st.cache[descr]
}

// BulkLookupRule returns a copy of a cached rule by tracker, if present.
func (pf *Client) BulkLookupRule(tracker string) (*FirewallRule, bool) {
	st := pf.bulk()
	st.mu.Lock()
	defer st.mu.Unlock()
	for i := range st.rulesCache {
		if st.rulesCache[i].Tracker == tracker {
			r := st.rulesCache[i]
			return &r, true
		}
	}
	return nil, false
}

// --- graceful-exit flush -------------------------------------------------
// Terraform kills the provider after an apply; Serve returns just before the
// process exits. Queued bulk ops must hit the box in that window, or the
// apply silently loses them when the queue is below the size threshold.

var (
	bulkClients   []*Client
	bulkClientsMu sync.Mutex
)

func (pf *Client) registerBulk() {
	bulkClientsMu.Lock()
	defer bulkClientsMu.Unlock()
	bulkClients = append(bulkClients, pf)
}

// FlushAllBulk flushes every registered client that still has queued ops.
func FlushAllBulk(ctx context.Context, timeout time.Duration) {
	bulkClientsMu.Lock()
	cl := append([]*Client(nil), bulkClients...)
	bulkClientsMu.Unlock()
	cctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	for _, c := range cl {
		if c.PendingCount() == 0 {
			continue
		}
		_ = c.FlushPending(cctx)
	}
}
