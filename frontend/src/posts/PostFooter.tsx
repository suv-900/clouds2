import React, { useEffect, useState } from "react";
import CustomError from "../components/CustomError";


//need 1 more work
export default function PostFooter(props:{
    tokenFound:boolean,
    likes:number,
    postid:number
}){
     
    const[displayError,setDisplayError] = useState(false);
    const[likes,setLikes] = useState(props.likes);
    const[liked,setLiked] = useState(false);    
    const[disliked,setDisliked] = useState(false);   

    const tokenFound = props.tokenFound;

    const[token,setToken] = useState<string>();
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

    async function likePost(){
        if(!props.tokenFound || token === undefined){
            renderErrorMessage();
            return;
        }
        if(disliked){
            await removeDislike()
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/likepost/${props.postid}`,{
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
        if(!props.tokenFound || token === undefined){
            renderErrorMessage();
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/removelike/${props.postid}`,{
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
    
    async function dislikePost(){
        if(!props.tokenFound || token === undefined){
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
        const response = await fetch(`http://localhost:8000/dislikepost/${props.postid}`,{
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

        if(!props.tokenFound || token === undefined){
            setDisplayError(true);
            vanishErrorMessage();
            return;
        }

        const headers = {
            "Authorization":token
        }
        const response = await fetch(`http://localhost:8000/removedislike/${props.postid}`,{
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
        <div>
            <div className="post-likes">{likes} likes</div>
            <button className="like-button" onClick={()=>{
                if(liked){
                    removeLike();
                }else{
                    likePost();
                }
                }}>{liked?"liked":"like"}</button>
            <button className="like-button" onClick={()=>{
                if(disliked){
                    removeDislike();
                }else{
                    dislikePost();
                }
                }}>{disliked?"disliked":"dislike"}</button>

            <CustomError enable={displayError} message={"please login"} />

        </div>
    )
}