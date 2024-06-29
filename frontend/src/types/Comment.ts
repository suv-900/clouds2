export default class Comment{
    id:number;
    authorid:number;
    authorname:string;
    content:string;
    likes:number;
    createdAt:string; 
    constructor(
        id:number,
        authorid:number,
        authorname:string,
        content:string,
        likes:number,
        createdAt:string
    ){
        this.id = id;
        this.authorid = authorid;
        this.authorname = authorname;
        this.content = content;
        this.likes = likes;
        this.createdAt = createdAt;
    }
}
