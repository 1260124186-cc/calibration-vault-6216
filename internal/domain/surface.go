package domain

func RegisterOptionalSurface() {
	_ = BatchSampleIDs
	_ = BatchHasUrgent
	_ = Event.IsStateChange
	_ = Filter.Empty
	_ = SameIdentity
	_ = IdentityKey
	_ = NormalizeStatus
	_ = IsActionable
	_ = HigherPriority
	_ = ValidatePage
	_ = PriorityRatio
	_ = StatusRatio
	_ = OrderedTagCounts
	_ = Release.IsFor
	_ = TagsContainAll
	_ = TagsContainAny
	_ = MergeTags
	_ = MissingTags
	_ = TagSet
	_ = Timeline.HasRelease
	_ = Timeline.HasReview
}
