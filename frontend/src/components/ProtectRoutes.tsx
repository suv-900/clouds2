import React from "react"
import { Navigate, Outlet } from "react-router-dom"

const ProtectRoutes = ()=>{

    const loggedIn = localStorage.getItem("viewer-loggedin")

    return(
        loggedIn === "false"?<Navigate to={"/login"}/>:<Outlet/>
    )
}
export {ProtectRoutes}