package v1alpha1

func (s ResourceClass) IsRequired() bool {
	for _, entry := range s.Spec.Entries {
		if entry.Required {
			return true
		}
	}
	return false
}

func (section *PanelSection) Contains(rd *ResourceDescriptor) bool {
	for _, entry := range section.Entries {
		if entry.Type != nil &&
			entry.Type.Group == rd.Spec.Resource.Group &&
			entry.Type.Version == rd.Spec.Resource.Version &&
			entry.Type.Resource == rd.Spec.Resource.Name {
			return true
		}
	}
	return false
}
