package config

type Config struct {
	RegistryRewrites map[string]string
	EnableRewrite    bool
	DBFile           string
	ForceDigest      bool
}
