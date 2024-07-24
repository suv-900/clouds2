import React, { useContext, useEffect, useState } from "react";
import Comment from "../../types/Comment";
import PostComment from "./PostComment";
import { AuthContext } from "./PostViewer";

export default function PostComments(props:{
    comments:Comment[] ,
    postid:number 
}){ 
    const[commentsList,setCommentsList] = useState(props.comments)
    const[comment,setComment] = useState<string>();
    const[displayError,setDisplayError] = useState(false);
    const[error,setError] = useState("");
    const[fetching,setFetching] = useState(false);
    const[offset,setOffset] = useState(0);
    const[hasMore,setHasmore] = useState(true);

    const token = useContext(AuthContext)
    
    useEffect(()=>{
        window.addEventListener("scroll",handleScroll)    
        setOffset(prevValue=>prevValue+5)
        
        return ()=>window.removeEventListener("scroll",handleScroll)
    },[])

    function handleScroll(){
        if(window.innerHeight + document.documentElement.scrollTop === document.documentElement.offsetHeight && !fetching){
            setFetching(true);
        }
    }

    useEffect(()=>{
        if(!fetching) return;
        setOffset(prevValue=>prevValue+5)
        getComments()
    },[fetching])

    async function getComments(){
        if(!hasMore){
            setFetching(false);
            return;
        }

        if(token){

        const response = await fetch(`http://localhost:8000/comments/getcomments?limit=5&offset=${offset}&postid=${props.postid}`,{
            method:"GET",
        })
        const res = await response.json();
        if(res === null){
            setHasmore(false);
            setFetching(false);
            return;
        }
        let commentsarr = commentsList
        for(let i =0;i<res.length;i++){
            let b = res[i];
            const comment = new Comment(
                b.Comment_id,
                b.User_id,
                b.Username,
                b.Comment_content,
                b.Comment_likes,
                b.Createdat_str,
                b.Liked,
                b.Disliked,
                false
            )
            commentsarr.push(comment);
        }
        setTimeout(()=>{
            setFetching(false);
            setCommentsList(commentsarr);
        },3000)
        }else{
            const response = await fetch(`http://localhost:8000/comments/getcomments?limit=5&offset=${offset}&postid=${props.postid}`,{
            method:"GET",
        })
        const res = await response.json();
        if(res === null){
            setHasmore(false);
            setFetching(false);
            return;
        }
        let commentsarr = commentsList
        for(let i =0;i<res.length;i++){
            let b = res[i];
            const comment = new Comment(
                b.Comment_id,
                b.User_id,
                b.Username,
                b.Comment_content,
                b.Comment_likes,
                b.Createdat_str,
                b.Liked,
                b.Disliked,
                false
            )
            commentsarr.push(comment);
        }
        setTimeout(()=>{
            setFetching(false);
            setCommentsList(commentsarr);
        },3000)
        }
    }


  
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
        if(token.length === 0){
            renderError("please login.")
            return;
        }
        if(comment === undefined || comment.length === 0 ){
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
            const curr = new Date();
            const timeStamp = curr.getHours()+":"+curr.getMinutes()+(curr.getHours() >= 12? ' pm':' am');
            
            const res = await response.json();
            let newComment = new Comment(
               res.Comment_id,
               res.User_id,
               res.Username,
               res.Comment_content,
               res.Comment_likes,
               timeStamp,
               false,
               false,
               true
            );
            setCommentsList([newComment,...commentsList]) 
        }else if(response.status === 500){
            renderError("try again.");
        }else if(response.status === 401 || 400){
            renderError("please login");
        }
    }

   
    return(
        <div >
            <div>
            {token.length !== 0 && props.postid?
        
        <div className="comment-box">
            <label>Add comment:</label><br></br>
            <textarea
            placeholder="write a comment..."
                onChange={
                    (e)=>{
                        setComment(e.target.value)
                }}
                /><br></br>
            <div className="cmb-buttondiv">
            <button className="post-like-button"onClick={()=>{addComment()}}>submit</button>
            </div>
        </div>
        :<></>}
        <div>{displayError?error:""}</div>
        
        </div>
            <div className="comments-section">
                <div>
                    <div className="comment-title">Comments: {commentsList.length}</div>
                    {commentsList.map((comment)=><PostComment comment={comment}/>)}
                </div>
                {fetching && 
                <div className="scroll-loader-div">
                    <span className="scroll-loader"></span>
                </div>
                }
            </div>
        </div>
    )
}