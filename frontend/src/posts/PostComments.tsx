import React, { useEffect } from "react";
import Comment from "../types/Comment";
import CommentBox from "./CommentBox";
import PostComment from "./PostComment";

export default function PostComments(props:{
    comments:Comment[] ,
}){ 

    return(
        <div >
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