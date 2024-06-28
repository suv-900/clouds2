import React, { useEffect, useState } from "react";
import Comment from "../types/Comment";
import CommentBox from "./CommentBox";
import PostComment from "./PostComment";
import { LinkedList } from "../utils/LinkedList";

export default function PostComments(props:{
    comments:Comment[] ,
    postid:number | undefined
}){ 

    const[list,setList] = useState<LinkedList<Comment>>();
    const[tokenFound,setTokenFound] = useState(false);
    const[token,setToken] = useState<string>();
    const[comment,setComment] = useState<string>();
    const[displayError,setDisplayError] = useState(false);
    const[error,setError] = useState("");
   
    useEffect(()=>{
        const token = localStorage.getItem("token");
        if(token !== null){
            setToken(token);
            setTokenFound(true);
        }
        let commentsList = new LinkedList<Comment>();

        for(const comment of props.comments){
            commentsList.push(comment);
        }
        console.log(commentsList);
        setList(commentsList);
    
    },[])
    
    function vanishErrorMessage(){
        setTimeout(()=>{
            setDisplayError(false);
            setError("");
        },3000)
    }
    function renderError(msg:string){
        setError(msg);
        setDisplayError(true);
        vanishErrorMessage();
    }
    async function addComment(){
        if(comment === undefined || comment.length === 0 || token === undefined){
            renderError("invalid.");
            return;
        }
        const headers={
            "Authorization":token
        }
        const body = JSON.stringify(comment);
        const response = await fetch(`http://localhost:8000/addcomment/${props.postid}`,{
            method:"POST",
            headers:headers,
            body:body,
        })
        if(response.ok){
            const res = await response.json();
            let newComment = new Comment(
               res.Comment_id,
               res.User_id,
               res.Username,
               res.Comment_content,
               res.Comment_likes 
            );
            console.log(newComment)
            if(list != undefined){
                list.push_front(newComment);
                setList(list);
            }else{
                let commentsList = new LinkedList<Comment>();
                commentsList.push(newComment);
                setList(commentsList);
            }
        }else if(response.status === 500){
            renderError("try again.");
        }else if(response.status === 401 || 400){
            renderError("please login");
        }
    }

   
    return(
        <div >
            <div>
            {tokenFound && props.postid?
        <div>
            <label>add comment</label>
            <input type="text" onChange={(e)=>{setComment(e.target.value)}}/>
            <button onClick={()=>{addComment()}}>submit</button>
        </div>:<></>}
        <div>{displayError?error:""}</div>
        
        </div>
           {props.comments ?
            <div className="comments-section">
                <div className="comment-title">Comments: {props.comments.length}</div>
                {props.comments.map((comment)=><PostComment 
                id={comment.id}
                content={comment.content}
                authorid={comment.authorid}
                authorname={comment.authorname}
                likes={comment.likes}
                />)} 
            </div>
            :<></>    
            } 
        </div>
    )
}