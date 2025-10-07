package internal

type RosterCreatorConfig struct {
	NumGK             int `mapstructure:"num_gk"`
	NumDF             int `mapstructure:"num_df"`
	NumDM             int `mapstructure:"num_dm"`
	NumMF             int `mapstructure:"num_mf"`
	NumAM             int `mapstructure:"num_am"`
	NumFW             int `mapstructure:"num_fw"`
	AvgStamina        int `mapstructure:"avg_stamina"`
	AvgAggression     int `mapstructure:"avg_aggression"`
	AvgMainSkill      int `mapstructure:"avg_main_skill"`
	AvgMidSkill       int `mapstructure:"avg_mid_skill"`
	AvgSecondarySkill int `mapstructure:"avg_secondary_skill"`
}

type Config struct {
	RosterCreator RosterCreatorConfig `mapstructure:"roster_creator"`
}
