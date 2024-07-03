import React from "react";
import "../css/styles.css"

function Loading(props:{enable:boolean}){
    return(
        <>
        {props.enable?
            <div className="loading-div">
            <div className="classic-2"></div> 
            </div>
        :<></>}
        </>
    )
}

export default Loading;