package data

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Comment struct {
	ID primitive.ObjectID

	PostID   primitive.ObjectID
	UserID   primitive.ObjectID
	Username string

	Content string
	Likes   int
}

type CommentClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func (cc CommentClient) AddComment(c context.Context, comment *Comment) (int64, error) {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	var commentID int64
	res, err := cc.db.Collection("comments").InsertOne(cx, comment)
	if err != nil {
		log.Error(err)
		return commentID, ErrInternalServerError
	}

	commentID, ok := res.InsertedID.(int64)
	if !ok {
		s := fmt.Sprintf("cannot convert %s interface{} type to int", res.InsertedID)
		log.Error(s)
		return commentID, ErrInternalServerError
	}

	return commentID, nil
}

func (cc CommentClient) DeleteComment(c context.Context, commentID primitive.ObjectID) error {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: commentID}}
	res, err := cc.db.Collection("comments").DeleteOne(cx, filter)

	if err != nil {
		log.Error(err)
		return ErrInternalServerError
	}

	if res.DeletedCount == 0 {
		return ErrRecordNotFound
	}

	return nil

}

func (cc CommentClient) UpdateComment(c context.Context, update bson.D, commentID primitive.ObjectID) error {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: commentID}}

	res, err := cc.db.Collection("comments").UpdateOne(cx, filter, update)

	if err != nil {
		log.Error(err)
		return ErrInternalServerError
	}

	if res.MatchedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (cc CommentClient) IncrementLike(c context.Context, commentID primitive.ObjectID) error {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: commentID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "likes", Value: 1}}}}

	res, err := cc.db.Collection("comments").UpdateOne(cx, filter, update)

	if err != nil {
		log.Error(err)
		return ErrInternalServerError
	}

	if res.MatchedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (cc CommentClient) DecrementLike(c context.Context, commentID primitive.ObjectID) error {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: commentID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "likes", Value: -1}}}}

	res, err := cc.db.Collection("comments").UpdateOne(cx, filter, update)

	if err != nil {
		log.Error(err)
		return ErrInternalServerError
	}

	if res.MatchedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}
