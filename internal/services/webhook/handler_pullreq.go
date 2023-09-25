// Copyright 2023 Harness, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webhook

import (
	"context"

	"github.com/harness/gitness/events"
	pullreqevents "github.com/harness/gitness/internal/events/pullreq"
	"github.com/harness/gitness/types"
	"github.com/harness/gitness/types/enum"
)

const (
	// gitReferenceNamePrefixBranch is the prefix of references of type branch.
	gitReferenceNamePrefixBranch = "refs/heads/"
)

// PullReqCreatedPayload describes the body of the pullreq created trigger.
// TODO: move in separate package for small import?
type PullReqCreatedPayload struct {
	BaseSegment
	PullReqSegment
	PullReqTargetReferenceSegment
	ReferenceSegment
	ReferenceDetailsSegment
}

// handleEventPullReqCreated handles created events for pull requests
// and triggers pullreq created webhooks for the source repo.
func (s *Service) handleEventPullReqCreated(ctx context.Context,
	event *events.Event[*pullreqevents.CreatedPayload]) error {
	return s.triggerForEventWithPullReq(ctx, enum.WebhookTriggerPullReqCreated,
		event.ID, event.Payload.PrincipalID, event.Payload.PullReqID,
		func(principal *types.Principal, pr *types.PullReq, targetRepo, sourceRepo *types.Repository) (any, error) {
			commitInfo, err := s.fetchCommitInfoForEvent(ctx, sourceRepo.GitUID, event.Payload.SourceSHA)
			if err != nil {
				return nil, err
			}
			targetRepoInfo := repositoryInfoFrom(targetRepo, s.urlProvider)
			sourceRepoInfo := repositoryInfoFrom(sourceRepo, s.urlProvider)

			return &PullReqCreatedPayload{
				BaseSegment: BaseSegment{
					Trigger:   enum.WebhookTriggerPullReqCreated,
					Repo:      targetRepoInfo,
					Principal: principalInfoFrom(principal),
				},
				PullReqSegment: PullReqSegment{
					PullReq: pullReqInfoFrom(pr),
				},
				PullReqTargetReferenceSegment: PullReqTargetReferenceSegment{
					TargetRef: ReferenceInfo{
						Name: gitReferenceNamePrefixBranch + pr.TargetBranch,
						Repo: targetRepoInfo,
					},
				},
				ReferenceSegment: ReferenceSegment{
					Ref: ReferenceInfo{
						Name: gitReferenceNamePrefixBranch + pr.SourceBranch,
						Repo: sourceRepoInfo,
					},
				},
				ReferenceDetailsSegment: ReferenceDetailsSegment{
					SHA:    event.Payload.SourceSHA,
					Commit: &commitInfo,
				},
			}, nil
		})
}

// PullReqReopenedPayload describes the body of the pullreq reopened trigger.
// Note: same as payload for created.
type PullReqReopenedPayload PullReqCreatedPayload

// handleEventPullReqReopened handles reopened events for pull requests
// and triggers pullreq reopened webhooks for the source repo.
func (s *Service) handleEventPullReqReopened(ctx context.Context,
	event *events.Event[*pullreqevents.ReopenedPayload]) error {
	return s.triggerForEventWithPullReq(ctx, enum.WebhookTriggerPullReqReopened,
		event.ID, event.Payload.PrincipalID, event.Payload.PullReqID,
		func(principal *types.Principal, pr *types.PullReq, targetRepo, sourceRepo *types.Repository) (any, error) {
			commitInfo, err := s.fetchCommitInfoForEvent(ctx, sourceRepo.GitUID, event.Payload.SourceSHA)
			if err != nil {
				return nil, err
			}
			targetRepoInfo := repositoryInfoFrom(targetRepo, s.urlProvider)
			sourceRepoInfo := repositoryInfoFrom(sourceRepo, s.urlProvider)

			return &PullReqReopenedPayload{
				BaseSegment: BaseSegment{
					Trigger:   enum.WebhookTriggerPullReqReopened,
					Repo:      targetRepoInfo,
					Principal: principalInfoFrom(principal),
				},
				PullReqSegment: PullReqSegment{
					PullReq: pullReqInfoFrom(pr),
				},
				PullReqTargetReferenceSegment: PullReqTargetReferenceSegment{
					TargetRef: ReferenceInfo{
						Name: gitReferenceNamePrefixBranch + pr.TargetBranch,
						Repo: targetRepoInfo,
					},
				},
				ReferenceSegment: ReferenceSegment{
					Ref: ReferenceInfo{
						Name: gitReferenceNamePrefixBranch + pr.SourceBranch,
						Repo: sourceRepoInfo,
					},
				},
				ReferenceDetailsSegment: ReferenceDetailsSegment{
					SHA:    event.Payload.SourceSHA,
					Commit: &commitInfo,
				},
			}, nil
		})
}

// PullReqBranchUpdatedPayload describes the body of the pullreq branch updated trigger.
// TODO: move in separate package for small import?
type PullReqBranchUpdatedPayload struct {
	BaseSegment
	PullReqSegment
	PullReqTargetReferenceSegment
	ReferenceSegment
	ReferenceDetailsSegment
	ReferenceUpdateSegment
}

// handleEventPullReqBranchUpdated handles branch updated events for pull requests
// and triggers pullreq branch updated webhooks for the source repo.
func (s *Service) handleEventPullReqBranchUpdated(ctx context.Context,
	event *events.Event[*pullreqevents.BranchUpdatedPayload]) error {
	return s.triggerForEventWithPullReq(ctx, enum.WebhookTriggerPullReqBranchUpdated,
		event.ID, event.Payload.PrincipalID, event.Payload.PullReqID,
		func(principal *types.Principal, pr *types.PullReq, targetRepo, sourceRepo *types.Repository) (any, error) {
			commitInfo, err := s.fetchCommitInfoForEvent(ctx, sourceRepo.GitUID, event.Payload.NewSHA)
			if err != nil {
				return nil, err
			}
			targetRepoInfo := repositoryInfoFrom(targetRepo, s.urlProvider)
			sourceRepoInfo := repositoryInfoFrom(sourceRepo, s.urlProvider)

			return &PullReqBranchUpdatedPayload{
				BaseSegment: BaseSegment{
					Trigger:   enum.WebhookTriggerPullReqBranchUpdated,
					Repo:      targetRepoInfo,
					Principal: principalInfoFrom(principal),
				},
				PullReqSegment: PullReqSegment{
					PullReq: pullReqInfoFrom(pr),
				},
				PullReqTargetReferenceSegment: PullReqTargetReferenceSegment{
					TargetRef: ReferenceInfo{
						Name: gitReferenceNamePrefixBranch + pr.TargetBranch,
						Repo: targetRepoInfo,
					},
				},
				ReferenceSegment: ReferenceSegment{
					Ref: ReferenceInfo{
						Name: gitReferenceNamePrefixBranch + pr.SourceBranch,
						Repo: sourceRepoInfo,
					},
				},
				ReferenceDetailsSegment: ReferenceDetailsSegment{
					SHA:    event.Payload.NewSHA,
					Commit: &commitInfo,
				},
				ReferenceUpdateSegment: ReferenceUpdateSegment{
					OldSHA: event.Payload.OldSHA,
					Forced: event.Payload.Forced,
				},
			}, nil
		})
}