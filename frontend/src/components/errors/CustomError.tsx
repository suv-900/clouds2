import React from "react";

export default function(props:{
    enable:boolean,
    message:string
}){
    return(
        <>
        {props.enable?
        <div className="error-message">{props.message}</div>
        :<></>}
        </>
    )
}