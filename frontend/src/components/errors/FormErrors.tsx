import React from "react"

function FormError(props:{
    render:boolean,
    msg:string | undefined
}){
    return(
        <div>
            {props.render && props.msg?<div>{props.msg}</div>:<></>}
        </div>
    )
}
export {FormError};