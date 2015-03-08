package onedrive

import "time"

// The HashesFacet groups different types of hashes into a single structure, for
// an item on OneDrive.
// See: http://onedrive.github.io/facets/hashes_facet.htm
type HashesFacet struct {
	Sha1Hash  string `json:"sha1Hash"`
	Crc32Hash string `json:"crc32Hash"`
}

// The FileFacet groups file-related data on OneDrive into a single structure.
// It is available on the file property of Item resources that represent files.
// See: http://onedrive.github.io/facets/file_facet.htm
type FileFacet struct {
	MimeType string       `json:"mimeType"`
	Hashes   *HashesFacet `json:"hashes"`
}

// The FolderFacet groups folder-related data on OneDrive into a single
// structure. It is available on the folder property of Item resources that
// represent folders.
// See: http://onedrive.github.io/facets/folder_facet.htm
type FolderFacet struct {
	ChildCount int64 `json:"childCount"`
}

// The ImageFacet groups image-related data on OneDrive into a single structure.
// See: http://onedrive.github.io/facets/image_facet.htm
type ImageFacet struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// The PhotoFacet groups photo-related data on OneDrive, for example, EXIF
// metadata, into a single structure.
// See: http://onedrive.github.io/facets/photo_facet.htm
type PhotoFacet struct {
	TakenDateTime       time.Time `json:"takenDateTime"`
	CameraMake          string    `json:"cameraMake"`
	CameraModel         string    `json:"cameraModel"`
	FNumber             float64   `json:"fNumber"`
	ExposureDenominator float64   `json:"exposureDenominator"`
	ExposureNumerator   float64   `json:"exposureNumerator"`
	FocalLength         float64   `json:"focalLength"`
	ISO                 float64   `json:"iso"`
}

// The AudioFacet groups audio-related data on OneDrive into a single structure.
// It is available on the audio property of Item resources that have associated audio.
// See: http://onedrive.github.io/facets/audio_facet.htm
type AudioFacet struct {
	Album             string  `json:"album"`
	AlbumArtist       string  `json:"albumArtist"`
	Artist            string  `json:"artist"`
	Bitrate           string  `json:"bitrate"`
	Composers         string  `json:"composers"`
	Copyright         string  `json:"copyright"`
	Disc              float64 `json:"disc"`
	DiscCount         float64 `json:"discCount"`
	Duration          float64 `json:"duration"`
	Genre             string  `json:"genre"`
	HasDRM            bool    `json:"hasDrm"`
	IsVariableBitrate bool    `json:"isVariableBitrate"`
	Title             string  `json:"title"`
	Track             float64 `json:"track"`
	TrackCount        float64 `json:"trackCount"`
	Year              float64 `json:"year"`
}

// The VideoFacet groups video-related data on OneDrive into a single complex type.
// See: http://onedrive.github.io/facets/video_facet.htm
type VideoFacet struct {
	Bitrate  string  `json:"bitrate"`
	Duration float64 `json:"duration"`
	Height   float64 `json:"height"`
	Width    float64 `json:"width"`
}

// The LocationFacet groups geographic location-related data on OneDrive into a
// single structure.
// See: http://onedrive.github.io/facets/location_facet.htm
type LocationFacet struct {
	Altitude  float64 `json:"altitude"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// The DeletedFacet indicates that the item on OneDrive has been deleted. In
// this version of the API, the presence (non-null) of the facet value indicates
// that the file was deleted. A null (or missing) value indicates that the file
//is not deleted.
// See: http://onedrive.github.io/facets/deleted_facet.htm
type DeletedFacet struct{}

// The SpecialFolder facet provides information about how a folder on OneDrive
// can be accessed via the special folders collection.
// See: http://onedrive.github.io/facets/jumpinfo_facet.htm
type SpecialFolder struct {
	Name string `json:"name"`
}

// The SharingLink type groups sharing link-related data on OneDrive into a
// single structure.
// See: http://onedrive.github.io/facets/sharinglink_facet.htm
type SharingLink struct {
	Token       string    `json:"token"`
	WebURL      string    `json:"webUrl"`
	Type        string    `json:"type"`
	Application *Identity `json:"application"`
}
