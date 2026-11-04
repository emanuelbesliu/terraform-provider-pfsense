# Vendored provider (THE build we run)

Source tree = fork of https://github.com/emanuelbesliu/terraform-provider-pfsense at release f16c920 (v0.58.0), plus our patches:
1. 35572af fix(firewall_rule): delete by GUI list position + CSRF + if=
2. 87c0a71 fix(aliases): tolerate write_config race in delete verification
3. 35de7ef feat(firewall_rule): floating/direction/quick + FloatingRules-tab delete
4. bulk apply mode (env PFSenseBulkApply=1, default OFF) + graceful-exit flush + proto-key alias + GUI-faithful empty fields (prov-src2@eece6c0, synced 2026-10-27)
5. bulk v2 (2026-11-03): flush returns the trackers it assigned (order-based, race-proof placeholder resolution), dead-placeholder self-heal (desc-absent + empty queue => state remove => clean recreate), log/disabled flags now honor Go's '1' values (php string-compare bug), deterministic state across tofu processes. Full lifecycle proven live on lab box: create76/update/refresh No-changes/bulk-delete-after-order-drift.
6. pfsense_lagg resource (2026-11-04, patch #6): first-class LAGG/LACP - pfSense stores laggs as an unkeyed list with a per-entry 'laggif' name and 'proto' field (verified against interfaces_lagg_edit.php + interface_lagg_configure on pfSense 2.9); name assigned server-side, apply=false default keeps remote changes config-plane only (self-cut protection), delete mirrors GUI (kernel destroy guarded by function_exists). Acceptance test TestAccLaggResource_basic (create+import+update+destroy) PASSES against virtual pfSense 2.9 lab VM: export TF_ACC=1 TF_ACC_TERRAFORM_PATH=/path/to/tofu TF_ACC_PROVIDER_HOST=registry.opentofu.org TF_ACC_PROVIDER_NAME=pfsense + TF_PFSENSE_URL/USERNAME/PASSWORD, go test ./internal/provider/ -run TestAccLagg.

Build from this vendor tree: `cd vendor/pfsense-provider && go build -o terraform-provider-pfsense_v0.58.1 .` (releases keep the v0.58.1 filename for the dev_override).
Binary in lab: terraform-provider-pfsense_v0.58.1, sha256 0c12a88470e7  # 2026-11-04 patch #7: pfsense_dhcp_backend + kea_custom_configuration (was bb276b6dd2d3)b7e4f83cbea2b3a2dcd2af2878b07cbbc0078b87e41aeaa2449b (identical build of this tree)
Upstream PRs pending (see docs/PROVIDER-PATCHES.md); internal mirror branch also pushed: SEKURE-960_provider_fork_fixes.

## 2026-10-23 acc-tests
Vendor copy = prov-src2@feat/rule-floating (delete+floating acceptance tests included).
The two skip-tagged floating-repros document real update bugs (HTTP 500) measured on the live box; upstream tracking in PR #91. Binary in ~/bin/devfixdir is source-identical for runtime behaviour (tests only touched _test files + README).
