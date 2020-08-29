package main

func (j *Jailguard) ListConfig() error {
	j.config.Print(j.cli.GetStdout(), "")
	return nil
}

func (j *Jailguard) SetConfigValue(k string, v string) error {
	err := j.config.Set(k, v)
	if err != nil {
		return err
	}
	err = j.config.Save()
	return err
}

func (j *Jailguard) ShowConfigValue(k string) error {
	j.config.Print(j.cli.GetStdout(), k)
	return nil
}
