export default class Comment{
    id:number;
    authorid:number;
    authorname:string;
    content:string;
    likes:number;
    createdAt:string; 
    userLiked:boolean;
    userDisliked:boolean
    newlyAddedComment:boolean;
    constructor(
        id:number,
        authorid:number,
        authorname:string,
        content:string,
        likes:number,
        createdAt:string,
        userLiked:boolean,
        userDisliked:boolean,
        newlyAddedComment:boolean
    ){
        this.id = id;
        this.authorid = authorid;
        this.authorname = authorname;
        this.content = content;
        this.likes = likes;
        this.createdAt = createdAt;
        this.userLiked = userLiked;
        this.userDisliked = userDisliked;
        this.newlyAddedComment = newlyAddedComment;
    }
}
