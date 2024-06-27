export default class Post{
    id:number;
    title:string;
    content:string;
    authorname:string;
    authorid:number;
    likes:number;
    createdat:string;

    constructor(
        id:number,
        title:string,
        content:string,
        authorname:string,
        authorid:number,
        likes:number,
        createdat:string
    ){
        this.id = id;
        this.title = title;
        this.content = content;
        this.authorname = authorname;
        this.authorid = authorid;
        this.likes = likes;
        this.createdat = createdat;
    }
    
}
