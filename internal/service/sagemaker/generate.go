// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

//go:generate go run ../../generate/tags/main.go -ListTags -ListTagsOp=ListTags -ListTagsOpPaginated -ServiceTagsSlice -TagOp=AddTags -UntagOp=DeleteTags -UpdateTags
//go:generate go run ../../generate/servicepackage/main.go
//go:generate go run ../../generate/listpages/main.go -ListOps=ListHubs
//go:generate go run ../../generate/identitytests/main.go
// ONLY generate directives and package declaration! Do not add anything else to this file.

package sagemaker
