//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/apps/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppConditions) DeepCopyInto(out *MyAppConditions) {
	*out = *in
	if in.PodInfoConditions != nil {
		in, out := &in.PodInfoConditions, &out.PodInfoConditions
		*out = make([]v1.DeploymentCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.RedisConditions != nil {
		in, out := &in.RedisConditions, &out.RedisConditions
		*out = make([]v1.DeploymentCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppConditions.
func (in *MyAppConditions) DeepCopy() *MyAppConditions {
	if in == nil {
		return nil
	}
	out := new(MyAppConditions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppImage) DeepCopyInto(out *MyAppImage) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppImage.
func (in *MyAppImage) DeepCopy() *MyAppImage {
	if in == nil {
		return nil
	}
	out := new(MyAppImage)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppRedisConfig) DeepCopyInto(out *MyAppRedisConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppRedisConfig.
func (in *MyAppRedisConfig) DeepCopy() *MyAppRedisConfig {
	if in == nil {
		return nil
	}
	out := new(MyAppRedisConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppResource) DeepCopyInto(out *MyAppResource) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppResource.
func (in *MyAppResource) DeepCopy() *MyAppResource {
	if in == nil {
		return nil
	}
	out := new(MyAppResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MyAppResource) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppResourceList) DeepCopyInto(out *MyAppResourceList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]MyAppResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppResourceList.
func (in *MyAppResourceList) DeepCopy() *MyAppResourceList {
	if in == nil {
		return nil
	}
	out := new(MyAppResourceList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *MyAppResourceList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppResourceSpec) DeepCopyInto(out *MyAppResourceSpec) {
	*out = *in
	in.Resources.DeepCopyInto(&out.Resources)
	out.Image = in.Image
	out.UI = in.UI
	out.Redis = in.Redis
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppResourceSpec.
func (in *MyAppResourceSpec) DeepCopy() *MyAppResourceSpec {
	if in == nil {
		return nil
	}
	out := new(MyAppResourceSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppResourceStatus) DeepCopyInto(out *MyAppResourceStatus) {
	*out = *in
	in.Conditions.DeepCopyInto(&out.Conditions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppResourceStatus.
func (in *MyAppResourceStatus) DeepCopy() *MyAppResourceStatus {
	if in == nil {
		return nil
	}
	out := new(MyAppResourceStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MyAppUI) DeepCopyInto(out *MyAppUI) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MyAppUI.
func (in *MyAppUI) DeepCopy() *MyAppUI {
	if in == nil {
		return nil
	}
	out := new(MyAppUI)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Resources) DeepCopyInto(out *Resources) {
	*out = *in
	out.MemoryLimit = in.MemoryLimit.DeepCopy()
	out.CPURequest = in.CPURequest.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Resources.
func (in *Resources) DeepCopy() *Resources {
	if in == nil {
		return nil
	}
	out := new(Resources)
	in.DeepCopyInto(out)
	return out
}
