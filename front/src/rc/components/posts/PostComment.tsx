import React, { useContext, useEffect, useState } from "react";
import CustomError from "../errors/CustomError";
import { AuthContext } from "./PostViewer";
import Comment from "../../types/Comment";
import CommentTimeStamp from "./NewlyAddedCommentTimeStamp";
import NewlyAddedCommentTimeStamp from "./NewlyAddedCommentTimeStamp";
import { BiSolidDislike, BiSolidLike } from "react-icons/bi";
import { getTime } from "../../utils/utils";

export default function PostComment(props:{
   comment:Comment 
}){

    const{comment} = props
    const token = useContext(AuthContext)
    const[timeStatus,setTimeStatus] = useState(getTime(comment.createdAt.substring(0,10)))
    const[likes,setLikes] = useState(comment.likes);
    const[liked,setLiked] = useState(comment.userLiked);    
    const[disliked,setDisliked] = useState(comment.userDisliked);    
    const[displayError,setDisplayError] = useState(false);
    
   
    function vanishErrorMessage(){
        setTimeout(()=>{
            setDisplayError(false);
        },3000)
    }

    function renderErrorMessage(){
        setDisplayError(true);
        vanishErrorMessage();
    }

    async function likeComment(){
        if(token.length === 0){
            renderErrorMessage();
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/likecomment/${comment.id}`,{
            method:"POST",
            headers:headers
        })
        if(response.ok){
            setLiked(true);
            setLikes(likes+1);

        }else if(response.status === 400){
            renderErrorMessage();
        }else{
            
        }
    }
    async function removeLike(){
        if(token.length === 0){
            renderErrorMessage();
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/removecommentlike/${comment.id}`,{
            method:"POST",
            headers:headers
        })
        if(response.ok){
            setLiked(false);
            setLikes(likes-1);
            console.log("removelike"+likes)
        }else if(response.status === 400){
            renderErrorMessage();
        }else{
            
        }
    }
    
    async function dislikeComment(){
        if(token.length === 0){
            renderErrorMessage();
            return;
        }

        if(liked){
            await removeLike()
        }
        
        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/dislikecomment/${comment.id}`,{
            method:"POST",
            headers:headers
        })
        if(response.ok){
            setDisliked(true);
            setLikes(likes-1);
            console.log("disliked"+likes)
        }else if(response.status === 400){
            renderErrorMessage();
        }else{
            
        }
    }
    async function removeDislike(){
        if(token.length === 0){
            renderErrorMessage();
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/removedislike/${comment.id}`,{
            method:"POST",
            headers:headers
        })
        if(response.ok){
            setDisliked(false);
            setLikes(likes+1);
        }else if(response.status === 400){
            renderErrorMessage();
        }else{
            
        }
    }

    const likebuttonstyles = {
        "color":liked?"#3bd139":""
    }
    const dislikebuttonstyles ={
        "color":disliked?"rgb(211, 38, 38)":""
    }
    return(
        <div id={comment.id.toString()}  className="comment">
            <a href={`http://localhost:3000/users/${comment.authorname}`} className="comment-authorname">{comment.authorname}</a>
            
            {comment.newlyAddedComment?<NewlyAddedCommentTimeStamp 
            newlyAddedComment={comment.newlyAddedComment}
            timeStamp={comment.createdAt}
            />:
            <time title={comment.createdAt.substring(0,17)} className="comment-createdat">{timeStatus}</time>
            }
            
            <div className="comment-content">{comment.content}</div>
            <div className="comment-likes">{likes} likes</div>
            <BiSolidLike size="20" style={likebuttonstyles} className="comment-like-button" onClick={()=>{
                if(liked){
                    removeLike();
                }else{
                    likeComment();
                }
                }}></BiSolidLike>
            <BiSolidDislike size="20" style={dislikebuttonstyles} className="comment-dislike-button" onClick={()=>{
                if(disliked){
                    removeDislike();
                }else{
                    dislikeComment();
                }
                }}>{disliked?"disliked":"dislike"}</BiSolidDislike>

            <CustomError enable={displayError} message={"please login"} />


        </div>
    )
}

