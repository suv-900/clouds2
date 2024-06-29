import React, { useEffect, useState } from "react";
import CustomError from "../components/CustomError";

export default function PostComment(props:{
    id:number,
    content:string,
    authorname:string,
    authorid:number,
    likes:number,
    createdAt:string
}){

    const[likes,setLikes] = useState(props.likes);
    const[liked,setLiked] = useState(false);    
    const[disliked,setDisliked] = useState(false);    
    const[token,setToken] = useState<string>();
    const[displayError,setDisplayError] = useState(false);
    
    useEffect(()=>{
        const token = localStorage.getItem("token");
        if(token){
            setToken(token);
        }
    },[])
    


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
        if(token === undefined){
            renderErrorMessage();
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/likecomment/${props.id}`,{
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
        if(token === undefined){
            renderErrorMessage();
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/removecommentlike/${props.id}`,{
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
        if(token === undefined){
            setDisplayError(true);
            vanishErrorMessage();
            return;
        }

        if(liked){
            await removeLike()
        }
        
        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/dislikecomment/${props.id}`,{
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
        if(token === undefined){
            setDisplayError(true);
            vanishErrorMessage();
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/removedislike/${props.id}`,{
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
    return(
        <div id={props.id.toString()}  className="comment">
            <a href="#">{props.authorname}</a>
            <div className="comment-content">{props.content}</div>
            <div className="comment-likes">{likes} likes</div>
            <div>{props.createdAt}</div>
            <button className="like-button" onClick={()=>{
                if(liked){
                    removeLike();
                }else{
                    likeComment();
                }
                }}>{liked?"liked":"like"}</button>
            <button className="like-button" onClick={()=>{
                if(disliked){
                    removeDislike();
                }else{
                    dislikeComment();
                }
                }}>{disliked?"disliked":"dislike"}</button>

            <CustomError enable={displayError} message={"please login"} />


        </div>
    )
}