package pfsense

import (
	"context"
	"fmt"
	"strings"
)

type packageResponse struct {
	Name             string `json:"name"`
	InstalledVersion string `json:"installed_version"`
	Description      string `json:"descr"`
	ShortName        string `json:"shortname"`
}

type Package struct {
	Name             string
	InstalledVersion string
	Description      string
	ShortName        string
}

func (p *Package) SetName(name string) error {
	if name == "" {
		return fmt.Errorf("%w, package name is required", ErrClientValidation)
	}

	p.Name = name

	return nil
}

type Packages []Package

func (pkgs Packages) GetByName(name string) (*Package, error) {
	for _, p := range pkgs {
		if p.Name == name || p.ShortName == name {
			return &p, nil
		}
	}

	return nil, fmt.Errorf("package %w with name '%s'", ErrNotFound, name)
}

func parsePackageResponse(resp packageResponse) Package {
	return Package{
		Name:             resp.Name,
		InstalledVersion: resp.InstalledVersion,
		Description:      resp.Description,
		ShortName:        resp.ShortName,
	}
}

func (pf *Client) getPackageFromPkgDB(ctx context.Context, name string) (*Package, error) {
	command := fmt.Sprintf(
		"$output = array();"+
			"$lines = array();"+
			"exec('/usr/local/sbin/pkg-static info -q \\'%s\\' 2>/dev/null', $lines, $rc);"+
			"if ($rc !== 0) {"+
			"$output['found'] = false;"+
			"} else {"+
			"$output['found'] = true;"+
			"$vlines = array();"+
			"exec('/usr/local/sbin/pkg-static query \\'%%v\\' \\'%s\\' 2>/dev/null', $vlines);"+
			"$output['version'] = count($vlines) > 0 ? $vlines[0] : '';"+
			"$dlines = array();"+
			"exec('/usr/local/sbin/pkg-static query \\'%%c\\' \\'%s\\' 2>/dev/null', $dlines);"+
			"$output['comment'] = count($dlines) > 0 ? $dlines[0] : '';"+
			"}"+
			"print(json_encode($output));",
		phpEscape(name),
		phpEscape(name),
		phpEscape(name),
	)

	var result struct {
		Found   bool   `json:"found"`
		Version string `json:"version"`
		Comment string `json:"comment"`
	}

	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w package from pkg database, %w", ErrGetOperationFailed, err)
	}

	if !result.Found {
		return nil, fmt.Errorf("package %w with name '%s' in pkg database", ErrNotFound, name)
	}

	return &Package{
		Name:             name,
		InstalledVersion: result.Version,
		Description:      result.Comment,
	}, nil
}

func (pf *Client) getPackages(ctx context.Context) (*Packages, error) {
	command := "require_once('pkg-utils.inc');" +
		"$pkgs = get_pkg_info('all', false, true);" +
		"$output = array();" +
		"foreach ($pkgs as $pkg) {" +
		"$item = array();" +
		"$item['name'] = $pkg['name'];" +
		"$item['installed_version'] = isset($pkg['installed_version']) ? $pkg['installed_version'] : '';" +
		"$item['descr'] = isset($pkg['descr']) ? $pkg['descr'] : '';" +
		"$item['shortname'] = isset($pkg['shortname']) ? $pkg['shortname'] : '';" +
		"array_push($output, $item);" +
		"};" +
		"print(json_encode($output));"

	var pkgResp []packageResponse
	if err := pf.executePHPCommand(ctx, command, &pkgResp); err != nil {
		return nil, err
	}

	pkgs := make(Packages, 0, len(pkgResp))
	for _, resp := range pkgResp {
		pkgs = append(pkgs, parsePackageResponse(resp))
	}

	return &pkgs, nil
}

func (pf *Client) GetPackages(ctx context.Context) (*Packages, error) {
	defer pf.read(&pf.mutexes.Package)()

	pkgs, err := pf.getPackages(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w packages, %w", ErrGetOperationFailed, err)
	}

	return pkgs, nil
}

func (pf *Client) GetPackage(ctx context.Context, name string) (*Package, error) {
	defer pf.read(&pf.mutexes.Package)()

	pkgs, err := pf.getPackages(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w packages, %w", ErrGetOperationFailed, err)
	}

	pkg, err := pkgs.GetByName(name)
	if err == nil {
		return pkg, nil
	}

	dbPkg, dbErr := pf.getPackageFromPkgDB(ctx, name)
	if dbErr != nil {
		return nil, fmt.Errorf("%w package, %w", ErrGetOperationFailed, err)
	}

	return dbPkg, nil
}

func (pf *Client) InstallPackage(ctx context.Context, name string) (*Package, error) {
	defer pf.write(&pf.mutexes.Package)()

	existingPkgs, err := pf.getPackages(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w packages for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingPkgs.GetByName(name); err == nil {
		return nil, fmt.Errorf("%w package, package '%s' is already installed", ErrCreateOperationFailed, name)
	}

	command := fmt.Sprintf(
		"require_once('pkg-utils.inc');"+
			"pkg_install('%s');"+
			"print(json_encode(true));",
		phpEscape(name),
	)

	var result bool
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w package '%s', %w", ErrCreateOperationFailed, name, err)
	}

	pkgs, err := pf.getPackages(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w packages after installing, %w", ErrGetOperationFailed, err)
	}

	pkg, err := pkgs.GetByName(name)
	if err != nil {
		return nil, fmt.Errorf("%w package after installing '%s', package not found in installed list — install may have failed silently, %w", ErrCreateOperationFailed, name, err)
	}

	return pkg, nil
}

func (pf *Client) InstallPackageFromURL(ctx context.Context, name string, packageURL string) (*Package, error) {
	defer pf.write(&pf.mutexes.Package)()

	existingPkgs, err := pf.getPackages(ctx)
	if err != nil {
		return nil, fmt.Errorf("%w packages for duplicate check, %w", ErrGetOperationFailed, err)
	}

	if _, err := existingPkgs.GetByName(name); err == nil {
		return nil, fmt.Errorf("%w package, package '%s' is already installed", ErrCreateOperationFailed, name)
	}

	if dbPkg, dbErr := pf.getPackageFromPkgDB(ctx, name); dbErr == nil {
		return dbPkg, nil
	}

	if !strings.HasPrefix(packageURL, "https://") && !strings.HasPrefix(packageURL, "http://") {
		return nil, fmt.Errorf("%w, package_url must start with http:// or https://", ErrClientValidation)
	}

	command := fmt.Sprintf(
		"$output = array();"+
			"$output['fetch'] = mwexec('/usr/local/sbin/pkg-static add -f \\'%s\\'');"+
			"$output['post_install'] = '';"+
			"$pkg_name = basename('%s');"+
			"$pkg_name = preg_replace('/\\.pkg$/', '', $pkg_name);"+
			"$post_install = '/usr/local/etc/rc.d/' . $pkg_name;"+
			"if (file_exists($post_install)) { $output['post_install'] = mwexec($post_install . ' start'); };"+
			"print(json_encode($output));",
		phpEscape(packageURL),
		phpEscape(packageURL),
	)

	var result map[string]interface{}
	if err := pf.executePHPCommand(ctx, command, &result); err != nil {
		return nil, fmt.Errorf("%w package '%s' from URL, %w", ErrCreateOperationFailed, name, err)
	}

	pkg, err := pf.getPackageFromPkgDB(ctx, name)
	if err != nil {
		pkgs, listErr := pf.getPackages(ctx)
		if listErr == nil {
			if p, findErr := pkgs.GetByName(name); findErr == nil {
				return p, nil
			}
		}

		return nil, fmt.Errorf("%w package after installing '%s' from URL, package not found in installed list — install may have failed, %w", ErrCreateOperationFailed, name, err)
	}

	return pkg, nil
}

func (pf *Client) DeletePackage(ctx context.Context, name string) error {
	defer pf.write(&pf.mutexes.Package)()

	existingPkgs, err := pf.getPackages(ctx)
	if err != nil {
		return fmt.Errorf("%w packages, %w", ErrGetOperationFailed, err)
	}

	inRepo := false
	if _, err := existingPkgs.GetByName(name); err == nil {
		inRepo = true
	} else {
		if _, dbErr := pf.getPackageFromPkgDB(ctx, name); dbErr != nil {
			return fmt.Errorf("%w package, %w", ErrGetOperationFailed, err)
		}
	}

	if inRepo {
		command := fmt.Sprintf(
			"require_once('pkg-utils.inc');"+
				"pkg_delete('%s');"+
				"print(json_encode(true));",
			phpEscape(name),
		)

		var result bool
		if err := pf.executePHPCommand(ctx, command, &result); err != nil {
			return fmt.Errorf("%w package '%s', %w", ErrDeleteOperationFailed, name, err)
		}
	} else {
		command := fmt.Sprintf(
			"$rc = mwexec('/usr/local/sbin/pkg-static delete -y \\'%s\\'');"+
				"print(json_encode(array('rc' => $rc)));",
			phpEscape(name),
		)

		var result map[string]interface{}
		if err := pf.executePHPCommand(ctx, command, &result); err != nil {
			return fmt.Errorf("%w package '%s', %w", ErrDeleteOperationFailed, name, err)
		}
	}

	pkgs, err := pf.getPackages(ctx)
	if err != nil {
		return fmt.Errorf("%w packages after deleting, %w", ErrGetOperationFailed, err)
	}

	if _, err := pkgs.GetByName(name); err == nil {
		return fmt.Errorf("%w package '%s', still installed after deletion", ErrDeleteOperationFailed, name)
	}

	if _, dbErr := pf.getPackageFromPkgDB(ctx, name); dbErr == nil {
		return fmt.Errorf("%w package '%s', still found in pkg database after deletion", ErrDeleteOperationFailed, name)
	}

	return nil
}
