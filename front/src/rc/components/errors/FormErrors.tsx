import React from "react"

function FormError(props:{
    enable:boolean,
    msg:string | undefined
}){
    return(
        <div>
            {props.enable && props.msg?<div 
            className="form-error-message">{props.msg}</div>:<></>}
        </div>
    )
}
export {FormError};