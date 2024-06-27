export default class Comment{
    id:number;
    authorid:number;
    authorname:string;
    content:string;
    likes:number;
    
    constructor(
        id:number,
        authorid:number,
        authorname:string,
        content:string,
        likes:number,
    ){
        this.id = id;
        this.authorid = authorid;
        this.authorname = authorname;
        this.content = content;
        this.likes = likes;
    }
}
