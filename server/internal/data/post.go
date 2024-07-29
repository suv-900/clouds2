package data

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Post struct {
	ID primitive.ObjectID

	Title      string
	Content    string
	AuthorID   uint64
	Authorname string
}

type PostClient struct {
	client *mongo.Client
	db     *mongo.Database
}

func (p PostClient) AddPost(c context.Context, post *Post) (int64, error) {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	var postID int64

	res, err := p.db.Collection("posts").InsertOne(cx, post)
	if err != nil {
		log.Error(err)
		return postID, ErrInternalServerError
	}

	postID, ok := res.InsertedID.(int64)
	if !ok {
		s := fmt.Sprintf("cannot convert %s interface{} type to int", res.InsertedID)
		log.Error(s)
		return postID, errors.New(s)
	}
	return postID, nil
}

func (p PostClient) GetPost(c context.Context, postID primitive.ObjectID) (Post, error) {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: postID}}

	var post Post
	err := p.db.Collection("posts").FindOne(cx, filter).Decode(&post)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return post, ErrRecordNotFound
		}
		return post, ErrInternalServerError
	}

	return post, nil
}

func (p PostClient) CheckPostTitleExists(c context.Context, post_title string) (bool, error) {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	filter := bson.D{{Key: "post_title", Value: post_title}}

	res := p.db.Collection("posts").FindOne(cx, filter)

	if res.Err() != nil {
		if res.Err() == mongo.ErrNoDocuments {
			return false, nil
		}
		return false, res.Err()
	}

	return true, nil
}

func (p PostClient) UpdatePost(c context.Context, update bson.D, postID primitive.ObjectID) error {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	opts := options.Update()
	filter := bson.D{{Key: "_id", Value: postID}}
	res, err := p.db.Collection("posts").UpdateOne(cx, filter, update, opts)

	if err != nil {
		log.Error(err)
		return ErrInternalServerError
	}

	if res.MatchedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (p PostClient) DeletePost(c context.Context, postID primitive.ObjectID) error {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	opts := options.Delete()
	filter := bson.D{{Key: "_id", Value: postID}}

	res, err := p.db.Collection("posts").DeleteOne(cx, filter, opts)

	if err != nil {
		log.Error(err)
		return ErrInternalServerError
	}

	if res.DeletedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// search by name
// capped collections and tailable cursor
func (p PostClient) GetPostsByAuthorID(c context.Context, authorID int64, offset int64) ([]Post, error) {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	var posts []Post

	filter := bson.D{{Key: "author_id", Value: authorID}}

	opts := options.Find().SetLimit(5).SetSkip(offset)

	cursor, err := p.db.Collection("posts").Find(cx, filter, opts)

	if err != nil {
		log.Error(err)
		return posts, ErrInternalServerError
	}

	if err := cursor.All(cx, posts); err != nil {
		return posts, ErrInternalServerError
	}

	return posts, nil
}

func (p PostClient) IncrementLike(c context.Context, postID primitive.ObjectID) error {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	opts := options.Update()
	filter := bson.D{{Key: "_id", Value: postID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "likes", Value: 1}}}}

	res, err := p.db.Collection("posts").UpdateOne(cx, filter, update, opts)

	if err != nil {
		log.Error(err)
		return ErrInternalServerError
	}
	if res.MatchedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}
func (p PostClient) DecrementLike(c context.Context, postID primitive.ObjectID) error {
	cx, cancel := context.WithTimeout(c, context_timeout)
	defer cancel()

	opts := options.Update()
	filter := bson.D{{Key: "_id", Value: postID}}
	update := bson.D{{Key: "$inc", Value: bson.D{{Key: "likes", Value: -1}}}}

	res, err := p.db.Collection("posts").UpdateOne(cx, filter, update, opts)

	if err != nil {
		log.Error(err)
		return ErrInternalServerError
	}
	if res.MatchedCount == 0 {
		return ErrRecordNotFound
	}

	return nil
}
