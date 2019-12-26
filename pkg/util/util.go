package util

import (
	corev1 "k8s.io/api/core/v1"
)

// Contains returns true if haystack contains needle
func Contains(haystack []string, needle string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}

// PodListContainsPod returns true if haystack contains needle
func PodListContainsPod(haystack corev1.PodList, needle corev1.Pod) bool {
	for _, p := range haystack.Items {
		if p.Name == needle.Name &&
			p.Namespace == needle.Namespace {
			return true
		}
	}
	return false
}
