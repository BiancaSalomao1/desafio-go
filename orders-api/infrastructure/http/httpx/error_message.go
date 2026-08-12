package httpx

func errorMessage(err error) string {

	if err == nil {
		return ""
	}

	return err.Error()
}
