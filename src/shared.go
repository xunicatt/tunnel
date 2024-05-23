package shared

const (
	TUNNEL_FILE = "tunnel.json"

	TUNNEL_VERSION   = "alpha~v.0.0.01a | go:debug"
	TUNNEL_MAJ_VER   = 0
	TUNNEL_MIN_VER   = 0
	TUNNEL_PATCH_VER = 1

	TUNNEL_DEF_PATH  = ".tunnel"
	TUNNEL_INIT_FILE = "init.json"
)

var (
	UserHomeDir string
	WorkingDir  string
	FileExists  bool
)
