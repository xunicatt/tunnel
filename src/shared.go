package shared

const (
	TUNNEL_FILE = "tunnel.json"

	TUNNEL_VERSION = "alpha~v.0.0.01a | go:debug | no-ver-increase\n" +
		"working= 'init' | 'install [downloads/zips/unzips/exapand-caches]'\n"
	TUNNEL_MAJ_VER   = 0
	TUNNEL_MIN_VER   = 0
	TUNNEL_PATCH_VER = 1

	TUNNEL_DEF_PATH   = ".tunnel"
	TUNNEL_CACHE_PATH = "cache"
	TUNNEL_PKG_PATH   = "pkg"
	TUNNEL_LIB_PATH   = "lib"
	TUNNEL_INIT_FILE  = "init.json"
)

var (
	UserHomeDir string
	WorkingDir  string
	FileExists  bool
)
