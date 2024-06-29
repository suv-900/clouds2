import React, { useEffect, useState } from "react";
import Post  from "../types/Post";
import Comment from "../types/Comment";
import Loading from "../components/Loading";
import PostComponent from "./PostComponent";

export default function PostViewer(){
    const[tokenFound,setTokenFound] = useState(false);
    const[loading,setLoading] = useState(true);
    const[clientError,setClientError] = useState(false);
    const[post,setPost] = useState<Post>()
    const[comments,setComments] = useState<Comment[]>([]);
    
    useEffect(()=>{
        const queryparams = new URLSearchParams(window.location.search);
        const i = queryparams.get("id"); 
        
        if(i === null){
            setClientError(true);
        }else{
            const id = parseInt(i); 
            getPost(id);
        } 
    },[])

    async function getPost(postid:number ){
        
        const token = localStorage.getItem("token");

        if(token !== null){
            const requestHeaders = {
                "Authorization":token
            }
            const response = await fetch(`http://localhost:8000/viewPostToken/${postid}`,{
                headers:requestHeaders,
                method:"GET"
            })
            if(response.ok){
                const res= await response.json();
                console.log(res);
                const k = res.Post;
                
                const post = new Post(
                    k.Post_id,
                    k.Post_title,
                    k.Post_content,
                    k.Author_name,
                    k.Author_id,
                    k.Post_likes,
                    res.PostLikedByUser,
                    res.PostDislikedByUser,
                    k.Createdat
                )
                setPost(post);

                const c = res.Comments;
                let commentsarr:Comment[] = [];
                for(let i =0;i<c.length;i++){
                    const b = c[i];
                    const comment = new Comment(
                        b.Comment_id,
                        b.User_id,
                        b.Username,
                        b.Comment_content,
                        b.Comment_likes,
                        b.CreatedAt,
                        b.Liked,
                        b.Disliked
                    )
                    commentsarr.push(comment);
                }
                setComments(commentsarr);
            }else if(response.status === 401){
                //component to delay and take back to login
            }else if(response.status === 500){
                // server error componenet
            }else{

            }
        }else{
            const response = await fetch(`http://localhost:8000/viewpost/${postid}`)
            if(response.ok){
                const res = await response.json();
                console.log(res);
                const k = res.Post;

                const post = new Post(
                    k.Post_id,
                    k.Post_title,
                    k.Post_content,
                    k.Author_name,
                    k.Author_id,
                    k.Post_likes,
                    false,
                    false,
                    k.Createdat
                )
                setPost(post);

                const c = res.Comments;
                let commentsarr:Comment[] = [];
                for(let i =0;i<c.length;i++){
                    const b = c[i];
                    const comment = new Comment(
                        b.Comment_id,
                        b.User_id,
                        b.Username,
                        b.Comment_content,
                        b.Comment_likes,
                        b.CreatedAt,
                        false,
                        false
                    )
                    commentsarr.push(comment);
                }
                setComments(commentsarr);

            }else if(response.status === 401){
                //component to delay and take back to login
            }else if(response.status === 500){
                // server error componenet
            }else{

            }
        }
        setTimeout(()=>{
            setLoading(false);
        },3000);
    }

    return(
        <div>
            <Loading enable={loading} />
            {!loading?
            <PostComponent 
            comments={comments}
            post={post}
            tokenFound={tokenFound}
            />
            :<></>}
        </div>
    )
}