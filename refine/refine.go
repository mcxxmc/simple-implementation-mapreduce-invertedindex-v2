package refine

func Refine(jobs []string) error {
	err := CreateSnapshots(jobs)
	if err != nil {
		return err
	}
	return nil
}
