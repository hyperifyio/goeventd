// Copyright (c) 2024. Heusala Group Oy <info@heusalagroup.fi>. All rights reserved.

package eventtrigger

type EventTrigger interface {
	TriggerEvent(serviceName string) bool
}
