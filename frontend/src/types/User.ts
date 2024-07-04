export default class User{
    userID:number
    username:string
    about:string
    createdat:string

    constructor(
        userID:number,
        username:string,
        about:string,
        createdat:string,
    ){
        this.userID = userID
        this.username = username
        this.about = about
        this.createdat = createdat
    }
}