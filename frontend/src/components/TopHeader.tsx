import React from "react";

export default function TopHeader(props:{
}){
    const blue = { color: 'blue' };
    const red = { color: 'red' };
    const orange = { color: 'orange' };
    const green = { color: 'green' };
    
    return(
        <div>
            <div className="top-header">
                <div className="header-text">
                    <span style={blue}>Welcome</span>
                    <span style={red}>To</span>
                    <span style={orange}>The</span>
                    <span style={green}>Internet!</span>
                </div>
                <div className="top-links">
                    <a href="http://localhost:3000/login">Login</a>
                    <a href="http://localhost:3000/register">Register</a>
                </div>
            </div>
            :<></>
        </div>
    )
}