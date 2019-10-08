/*

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

package controllers

import (
    "context"

    "github.com/go-logr/logr"
    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "golang.org/x/net/context"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "google.golang.org/api/youtube/v3"

    youtubev1alpha1 "pablokbs.com/m/api/v1alpha1"
)

// YouTubeVideoReconciler reconciles a YouTubeVideo object
type YouTubeVideoReconciler struct {
    client.Client
    Log logr.Logger
}

func ignoreNotFound(err error) error {
    if apierrs.IsNotFound(err) {
        return nil
    }
    return err
}

// +kubebuilder:rbac:groups=youtube.github.com/pablokbs/youtube-crd,resources=youtubevideoes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=youtube.github.com/pablokbs/youtube-crd,resources=youtubevideoes/status,verbs=get;update;patch

func (r *YouTubeVideoReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
    ctx := context.Background()
    log := r.Log.WithValues("youtubevideo", req.NamespacedName)

    // your logic here

    var myVideos youtubev1alpha1.YouTubeVideo
    if err := r.Get(ctx, req.NamespacedName, &myVideos); err != nil {
        log.Error(err, "unable to fetch Server")
        // we'll ignore not-found errors, since they can't be fixed by an immediate
        // requeue (we'll need to wait for a new notification), and we can get them
        // on deleted requests.
        return ctrl.Result{}, ignoreNotFound(err)
    }

    return ctrl.Result{}, nil
}

func (r *ServerReconciler) getVideos(s *youtubev1alpha1.YouTubeVideo) (*godo.KubernetesClusterCreateRequest, error) {

    var (
        method                  = flag.String("method", "list", "The API method to execute. (List is the only method that this sample currently supports.")
        channelId               = flag.String("channelId", "", "Retrieve playlists for this channel. Value is a YouTube channel ID.")
        hl                      = flag.String("hl", "", "Retrieve localized resource metadata for the specified application language.")
        maxResults              = flag.Int64("maxResults", 5, "The maximum number of playlist resources to include in the API response.")
        mine                    = flag.Bool("mine", false, "List playlists for authenticated user's channel. Default: false.")
        onBehalfOfContentOwner  = flag.String("onBehalfOfContentOwner", "", "Indicates that the request's auth credentials identify a user authorized to act on behalf of the specified content owner.")
        pageToken               = flag.String("pageToken", "", "Token that identifies a specific page in the result set that should be returned.")
        part                    = flag.String("part", "snippet", "Comma-separated list of playlist resource parts that API response will include.")
        playlistId              = flag.String("playlistId", "", "Retrieve information about this playlist.")
    )

    call := service.PlaylistItems.List(part)
    call = call.PlaylistId(playlistId)
    if pageToken != "" {
            call = call.PageToken(pageToken)
    }
    response, err := call.Do()
    handleError(err, "")
    return response
}

func (r *YouTubeVideoReconciler) SetupWithManager(mgr ctrl.Manager) error {
    return ctrl.NewControllerManagedBy(mgr).
        For(&youtubev1alpha1.YouTubeVideo{}).
        Complete(r)
}
