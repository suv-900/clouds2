import React, { useState } from "react";

export default function NewlyAddedCommentTimeStamp(props:{
    newlyAddedComment:boolean,
    timeStamp:string
}){
    const[timeStamp,setTimeStamp] = useState("Just now")

    setTimeout(()=>{
        setTimeStamp(props.timeStamp)
    },4000)

    return(
        <time title={timeStamp} className="comment-createdat">{timeStamp}</time>
    )
}