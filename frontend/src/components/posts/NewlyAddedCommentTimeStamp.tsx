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
        <div className="comment-createdat">{timeStamp}</div>
    )
}