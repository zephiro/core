package handle

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/tryanzu/core/board/legacy/model"
	"github.com/tryanzu/core/core/events"
	"github.com/tryanzu/core/deps"
	"github.com/tryanzu/core/modules/gaming"
	"github.com/tryanzu/core/modules/user"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
)

type VoteAPI struct {
	Gaming *gaming.Module `inject:""`
	User   *user.Module   `inject:""`
}

func (di VoteAPI) VoteComponent(c *gin.Context) {

	// Get the database interface from the DI
	database := deps.Container.Mgo()

	id := c.Params.ByName("id")

	if bson.IsObjectIdHex(id) == false {

		// Invalid request
		c.JSON(400, gin.H{"error": "Invalid request...", "status": 601})

		return
	}

	// Get the user
	user_id := c.MustGet("user_id")
	user_bson_id := bson.ObjectIdHex(user_id.(string))

	// Get the vote content
	var vote model.VoteForm

	if c.BindWith(&vote, binding.JSON) == nil {

		if vote.Direction == "up" || vote.Direction == "down" {

			// Check if component is valid
			component := vote.Component
			direction := vote.Direction
			valid := false

			for _, possible := range avaliable_components {

				if component == possible {

					valid = true
				}
			}

			if valid == true {

				// Add the vote itself to the votes collection
				var value int

				if direction == "up" {
					value = 1
				}

				if direction == "down" {
					value = -1
				}

				// Get the post using the slug
				id := bson.ObjectIdHex(id)
				collection := database.C("posts")

				var post model.Post
				err := collection.FindId(id).One(&post)

				if err != nil {

					// No guest can vote
					c.JSON(404, gin.H{"error": "No post found...", "status": 605})

					return
				} else {

					var add bytes.Buffer

					// Make the push string
					add.WriteString("components.")
					add.WriteString(component)
					add.WriteString(".votes.")
					add.WriteString(direction)

					inc := add.String()

					var already_voted model.Vote

					err = database.C("votes").Find(bson.M{"type": "component", "user_id": user_bson_id, "related_id": id, "nested_type": component}).One(&already_voted)

					if err == nil {

						var rem bytes.Buffer

						// Make the push string
						rem.WriteString("components.")
						rem.WriteString(component)
						rem.WriteString(".votes.")

						if (direction == "up" && already_voted.Value == 1) || (direction == "down" && already_voted.Value == -1) {

							rem.WriteString(direction)
							ctc := rem.String()
							change := bson.M{"$inc": bson.M{ctc: -1}}
							err = collection.Update(bson.M{"_id": post.Id}, change)

							if err != nil {

								panic(err)
							}

							err = database.C("votes").RemoveId(already_voted.Id)

							if err != nil {

								panic(err)
							}

							c.JSON(200, gin.H{"message": "okay", "status": 609})
							return

						} else if (direction == "up" && already_voted.Value == -1) || (direction == "down" && already_voted.Value == 1) {

							if direction == "up" {

								rem.WriteString("down")

							} else if direction == "down" {

								rem.WriteString("up")
							}

							ctc := rem.String()

							change := bson.M{"$inc": bson.M{ctc: -1}}

							err = collection.Update(bson.M{"_id": post.Id}, change)

							if err != nil {

								panic(err)
							}

							err = database.C("votes").RemoveId(already_voted.Id)

							if err != nil {

								panic(err)
							}
						}
					}

					change := bson.M{"$inc": bson.M{inc: 1}}
					err = collection.Update(bson.M{"_id": post.Id}, change)

					if err != nil {

						panic(err)
					}

					vote := &model.Vote{
						UserId:     user_bson_id,
						Type:       "component",
						NestedType: component,
						RelatedId:  id,
						Value:      value,
						Created:    time.Now(),
					}
					err = database.C("votes").Insert(vote)

					c.JSON(200, gin.H{"message": "okay", "status": 606})
					return
				}

			} else {

				// Component does not exists
				c.JSON(400, gin.H{"error": "Not authorized...", "status": 604})
			}
		}
	}

	c.JSON(401, gin.H{"error": "Couldnt create post, missing information...", "status": 205})
}

func (di VoteAPI) VoteComment(c *gin.Context) {

	// Get the database interface from the DI
	database := deps.Container.Mgo()

	id := c.Params.ByName("id")

	if bson.IsObjectIdHex(id) == false {

		// Invalid request
		c.JSON(400, gin.H{"error": "Invalid request...", "status": 601})

		return
	}

	// Get the user
	user_id := c.MustGet("user_id")
	user_bson_id := bson.ObjectIdHex(user_id.(string))

	// Get the vote content
	var vote model.VoteCommentForm

	if c.BindWith(&vote, binding.JSON) == nil {

		// Get the post using the id
		id := bson.ObjectIdHex(id)
		collection := database.C("posts")

		var post model.Post
		err := collection.FindId(id).One(&post)

		if err != nil {
			panic(err)
		}

		// Get the author of the vote
		usr, err := di.User.Get(user_bson_id)

		if err != nil {

			c.JSON(400, gin.H{"status": "error", "message": err.Error()})
			return
		}

		user_model := usr.Data()
		index := vote.Comment

		if _, err := strconv.Atoi(index); err == nil {

			var add bytes.Buffer
			var already_voted model.Vote
			var vote_value int

			err = database.C("votes").Find(bson.M{"type": "comment", "user_id": user_bson_id, "related_id": id, "nested_type": index}).One(&already_voted)

			comment_index, _ := strconv.Atoi(index)
			comment_ := post.Comments.Set[comment_index]

			if err == nil {

				// Cannot allow to change a vote once 15 minutes have passed
				if time.Since(already_voted.Created) > time.Minute*15 {

					c.JSON(400, gin.H{"message": "Cannot allow vote changes after 15 minutes.", "status": "error"})

					return
				}

				var rem bytes.Buffer

				// Remove the vote from the comment using $inc
				rem.WriteString("comments.set.")
				rem.WriteString(index)

				if already_voted.Value == 1 {

					rem.WriteString(".votes.up")
				}

				if already_voted.Value == -1 {

					rem.WriteString(".votes.down")
				}

				// Comment-To-Change
				ctc := rem.String()

				change := bson.M{"$inc": bson.M{ctc: -1}}
				err = collection.Update(bson.M{"_id": post.Id}, change)

				if err != nil {
					panic(err)
				}

				err = database.C("votes").RemoveId(already_voted.Id)

				if err != nil {
					panic(err)
				}

				// Return the gamification points
				if already_voted.Value == 1 {

					go func(usr *user.One, comment model.Comment) {

						di.Gaming.Get(usr).Tribute(1)

						author := di.Gaming.Get(comment.UserId)
						author.Coins(-1)

						if comment.Votes.Up-1 < 5 {
							author.Swords(-1)
						}

					}(usr, comment_)

					events.In <- events.RawEmit("post", post.Id.Hex(), map[string]interface{}{
						"fire":  "comment-upvote-remove",
						"index": comment_index,
					})

				} else if already_voted.Value == -1 {

					go func(usr *user.One, comment model.Comment) {

						di.Gaming.Get(usr).Shit(1)

						author := di.Gaming.Get(comment.UserId)
						author.Coins(1)

						if comment.Votes.Down-1 < 5 {
							author.Swords(1)
						}

					}(usr, comment_)

					events.In <- events.RawEmit("post", post.Id.Hex(), map[string]interface{}{
						"fire":  "comment-downvote-remove",
						"index": comment_index,
					})
				}

				c.JSON(200, gin.H{"status": "okay"})
				return
			}

			// Check if has enough tribute or shit to give
			if (vote.Direction == "up" && user_model.Gaming.Tribute < 1) || (vote.Direction == "down" && user_model.Gaming.Shit < 1) {
				c.JSON(400, gin.H{"message": "Dont have enough gaming points to do this.", "status": "error"})
				return
			}

			// Make the push string
			add.WriteString("comments.set.")
			add.WriteString(index)

			if vote.Direction == "up" {

				vote_value = 1
				add.WriteString(".votes.up")

				events.In <- events.RawEmit("post", post.Id.Hex(), map[string]interface{}{
					"fire":  "comment-upvote",
					"index": comment_index,
				})
			}

			if vote.Direction == "down" {
				vote_value = -1
				add.WriteString(".votes.down")

				events.In <- events.RawEmit("post", post.Id.Hex(), map[string]interface{}{
					"fire":  "comment-downvote",
					"index": comment_index,
				})
			}

			inc := add.String()

			change := bson.M{"$inc": bson.M{inc: 1}}
			err = collection.Update(bson.M{"_id": post.Id}, change)

			if err != nil {
				panic(err)
			}

			vote := &model.Vote{
				UserId:     user_bson_id,
				Type:       "comment",
				NestedType: index,
				RelatedId:  id,
				Value:      vote_value,
				Created:    time.Now(),
			}

			err = database.C("votes").Insert(vote)

			// Remove the spend of tribute or shit when giving the vote to the comment (only if comment's user is not the same as the vote's user)
			if comment_.UserId != user_bson_id {

				if vote_value == -1 {

					go func(usr *user.One, comment model.Comment) {

						di.Gaming.Get(usr).Shit(-1)

						author := di.Gaming.Get(comment.UserId)
						author.Coins(-1)

						if comment.Votes.Down <= 5 {
							author.Swords(-1)
						}

					}(usr, comment_)

				} else {

					go func(usr *user.One, comment model.Comment) {

						di.Gaming.Get(usr).Tribute(-1)

						author := di.Gaming.Get(comment.UserId)
						author.Coins(1)

						if comment.Votes.Up <= 5 {
							author.Swords(1)
						}

					}(usr, comment_)
				}
			}

			// Notify the author of the comment
			go func(comment model.Comment, token bson.ObjectId, post model.Post) {

				/*user_id := comment.UserId

				  // Get the comment like author
				  var user model.User

				  database.C("users").Find(bson.M{"_id": token.UserId}).One(&user)

				  if err == nil {

				      // Gravatar url
				      emailHash := gravatar.EmailHash(user.Email)
				      image := gravatar.GetAvatarURL("http", emailHash, "http://spartangeek.com/images/default-avatar.png", 80)

				      // Construct the notification message
				      title := fmt.Sprintf("A **%s** le gusta tu comentario.", user.UserName)
				      message := post.Title

				      // We are inside an isolated routine, so we dont need to worry about the processing cost
				      //notify(user_id, "like", post.Id, "/post/" + post.Slug, title, message, image.String())
				  }*/

			}(comment_, user_bson_id, post)

			c.JSON(200, gin.H{"status": "okay"})
			return
		}
	}

	c.JSON(401, gin.H{"error": "Couldnt vote, missing information...", "status": 608})
}

func (di VoteAPI) VotePost(c *gin.Context) {

	// Get the database interface from the DI
	database := deps.Container.Mgo()

	id := c.Params.ByName("id")

	if bson.IsObjectIdHex(id) == false {
		c.JSON(400, gin.H{"error": "Invalid request...", "status": 601})
		return
	}

	// Get the user
	user_id := c.MustGet("user_id")
	user_bson_id := bson.ObjectIdHex(user_id.(string))

	// Get the vote content
	var vote model.VotePostForm

	if c.BindWith(&vote, binding.JSON) == nil {

		// Get the post using the id
		id := bson.ObjectIdHex(id)
		collection := database.C("posts")

		var post model.Post
		err := collection.FindId(id).One(&post)

		if err != nil {
			panic(err)
		}

		// Get the author of the vote
		usr, err := di.User.Get(user_bson_id)
		if err != nil {
			c.JSON(400, gin.H{"status": "error", "message": err.Error()})
			return
		}

		user_model := usr.Data()

		var add bytes.Buffer
		var already_voted model.Vote
		var vote_value int

		err = database.C("votes").Find(bson.M{"type": "post", "user_id": user_bson_id, "related_id": id}).One(&already_voted)

		if err == nil {

			// Cannot allow to change a vote once 15 minutes have passed
			if time.Since(already_voted.Created) > time.Minute*15 {

				c.JSON(400, gin.H{"message": "Cannot allow vote changes after 15 minutes.", "status": "error"})

				return
			}

			var rem bytes.Buffer

			// Remove the vote from the comment using $inc
			rem.WriteString("votes.")

			if already_voted.Value == 1 {

				rem.WriteString("up")
			}

			if already_voted.Value == -1 {

				rem.WriteString("down")
			}

			// Comment-To-Change
			ctc := rem.String()

			change := bson.M{"$inc": bson.M{ctc: -1}}
			err = collection.Update(bson.M{"_id": post.Id}, change)

			if err != nil {
				panic(err)
			}

			err = database.C("votes").RemoveId(already_voted.Id)

			if err != nil {
				panic(err)
			}

			// Return the gamification points
			if already_voted.Value == 1 {

				go func(usr *user.One, post model.Post) {

					di.Gaming.Get(usr).Tribute(1)

					author := di.Gaming.Get(post.UserId)
					author.Coins(-2)

					if post.Votes.Up-1 < 10 {
						author.Swords(-1)
					}

				}(usr, post)

				events.In <- events.RawEmit("post", post.Id.Hex(), map[string]interface{}{
					"fire": "upvote-remove",
				})

			} else if already_voted.Value == -1 {

				go func(usr *user.One, post model.Post) {

					di.Gaming.Get(usr).Shit(1)

					author := di.Gaming.Get(post.UserId)
					author.Coins(1)

					if post.Votes.Down-1 < 10 {
						author.Swords(1)
					}

				}(usr, post)

				events.In <- events.RawEmit("post", post.Id.Hex(), map[string]interface{}{
					"fire": "downvote-remove",
				})
			}

			c.JSON(200, gin.H{"status": "okay"})
			return
		}

		// Check if has enough tribute or shit to give
		if (vote.Direction == "up" && user_model.Gaming.Tribute < 1) || (vote.Direction == "down" && user_model.Gaming.Shit < 1) {

			c.JSON(400, gin.H{"message": "Dont have enough gaming points to do this.", "status": "error"})

			return
		}

		// Make the push string
		add.WriteString("votes.")

		if vote.Direction == "up" {
			vote_value = 1
			add.WriteString("up")

			events.In <- events.RawEmit("post", post.Id.Hex(), map[string]interface{}{
				"fire": "upvote",
			})
		}

		if vote.Direction == "down" {
			vote_value = -1
			add.WriteString("down")

			events.In <- events.RawEmit("post", post.Id.Hex(), map[string]interface{}{
				"fire": "downvote",
			})
		}

		inc := add.String()

		change := bson.M{"$inc": bson.M{inc: 1}}
		err = collection.Update(bson.M{"_id": post.Id}, change)

		if err != nil {
			panic(err)
		}

		vote := &model.Vote{
			UserId:    user_bson_id,
			Type:      "post",
			RelatedId: id,
			Value:     vote_value,
			Created:   time.Now(),
		}

		err = database.C("votes").Insert(vote)

		// Remove the spend of tribute or shit when giving the vote to the comment (only if comment's user is not the same as the vote's user)
		if post.UserId != user_bson_id {

			if vote_value == -1 {

				go func(usr *user.One, post model.Post) {

					di.Gaming.Get(usr).Shit(-1)

					author := di.Gaming.Get(post.UserId)
					author.Coins(-1)

					if post.Votes.Down <= 10 {
						author.Swords(-1)
					}

				}(usr, post)

			} else {

				go func(usr *user.One, post model.Post) {

					di.Gaming.Get(usr).Tribute(-1)

					author := di.Gaming.Get(post.UserId)
					author.Coins(2)

					if post.Votes.Up <= 10 {
						author.Swords(1)
					}

				}(usr, post)
			}
		}

		c.JSON(200, gin.H{"status": "okay"})
		return

	}

	c.JSON(401, gin.H{"error": "Couldnt vote, missing information...", "status": 608})
}
