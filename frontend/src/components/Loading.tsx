import React from "react";
import "../css/styles.css"

function Loading(props:{enable:boolean}){
    return(
        <>
        {props.enable?
            <div className="classic-2"></div> 
        :<></>}
        </>
    )
}

export default Loading;