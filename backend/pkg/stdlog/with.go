package stdlog

func (i *impl) With(tags map[string]interface{}) Logger {
	if len(tags) == 0 {
		return i
	}

	newTags := normalizeTags(tags)
	logger := i.logger
	return &impl{
		cfg:    i.cfg,
		logger: logger.With(zapFields(newTags)...),
	}
}
