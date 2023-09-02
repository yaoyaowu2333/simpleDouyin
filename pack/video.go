package pack

import (
	"simpleDouyin/dao"
	"simpleDouyin/entity"
)

// AuthorIds
// 依据视频内容获取返回视频的作者id
func AuthorIds(videoModels []*dao.Video) []int64 {
	if videoModels != nil {
		var ids = make([]int64, 0, len(videoModels))
		for _, videoModel := range videoModels {
			ids = append(ids, videoModel.AuthorId)
		}
		return ids
	}
	return []int64{}
}

func Video(videoModel *dao.Video) *entity.Video {
	if videoModel != nil {
		return &entity.Video{
			Id:            videoModel.Id,
			Author:        entity.User{},
			PlayUrl:       videoModel.PlayUrl,
			CoverUrl:      videoModel.CoverUrl,
			FavoriteCount: videoModel.FavoriteCount,
			CommentCount:  videoModel.FavoriteCount,
			Title:         videoModel.Title,
		}
	}
	return nil
}

// Videos
// 将videoModels数据通过Videos进行处理，在拷贝的过程中对数据进行组装
func Videos(videoModels []*dao.Video) []*entity.Video {
	if videoModels != nil {
		var videos = make([]*entity.Video, 0, len(videoModels))
		for _, model := range videoModels {
			videos = append(videos, Video(model))
		}
		return videos
	}
	return nil
}

func VideoPtrs(videoPtrs []*entity.Video) []entity.Video {
	if videoPtrs != nil {
		var videos = make([]entity.Video, len(videoPtrs))
		for i, ptr := range videoPtrs {
			videos[i] = *ptr
		}
		return videos
	}
	return []entity.Video{}
}
