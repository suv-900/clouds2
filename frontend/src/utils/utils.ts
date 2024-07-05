function checkPassword(password:string):boolean{
    if(password.length === 0) return false;

    const len =(password.length >=8 && password.length <= 20)?true:false;
    const digit = /\d/.test(password);
    const specialChar = /[@#$%&]/.test(password);
    
    return len && digit && specialChar;
}
function checkEmail(email:string):boolean{
    if(email.length === 0) return false;

    return /@gmail.com/.test(email);
}

// function decodeTime(s:string):string{
//     const currentTime = Date() 
//     let date="";
//     let time="";
//     let i = 0;
//     while(s[i] !== 'T'){
//         date += s[i]
//     }
// }
function getTime(s:string):string{
    const curr = new Date();
    
    const curryear = curr.getFullYear()
    const currmonth = curr.getMonth()
    const currday = curr.getDay()

    const past = new Date(s.substring(0,17))
    const pastyear = past.getFullYear()
    const pastmonth = past.getMonth()
    const pastday = past.getDay()

    if((curryear - pastyear) > 0){
        const diff = curryear - pastyear
        if(diff === 1){
            return `a year ago`;
        }else{
            return `${curryear-pastyear} year ago`;
        }
    }else if((currmonth - pastmonth) > 0){
        const diff = currmonth - pastmonth
        if(diff === 1){
            return `a month ago`;
        }else{
            return `${currmonth - pastmonth} months ago`;
        }
    }else if((currday - pastday) > 0){
        const diff = currday - pastday
        if(diff === 1){
            return `a day ago`;
        }else{
            return `${currday - pastday} days ago`;
        }
    }else{
        return `today`
    }
}
const months = [
        "Jan", "Feb", "Mar", "Apr", "May", "Jun",
        "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"
    ];
function getTimeForPosts(s:string):string{
    const curr = new Date();
    
    const curryear = curr.getFullYear()
    const currmonth = curr.getMonth()
    const currday = curr.getDay()

    const past = new Date(s.substring(0,10))
    const pastyear = past.getFullYear()
    const pastmonth = past.getMonth()
    const pastday = past.getDay()
    console.log(s);
    console.log(past)
    console.log(pastyear+","+pastmonth+","+pastday);

    if((curryear - pastyear) > 0){
        const diff = curryear - pastyear
        return months[(currmonth - pastmonth)-1]+","+pastyear.toString() 
    }else if((currmonth - pastmonth) > 0){
        const diff = currmonth - pastmonth 
        return months[diff - 1]+","+pastday.toString() 
    }else if((currday - pastday) > 0){
        const diff = currday - pastday
        if(diff === 1){
            return `a day ago`;
        }else{
            const diff = currmonth - pastmonth 
            return months[diff - 1]+","+pastday.toString() 
        }
    }else{
        return `today`
    }
}
export {checkPassword,checkEmail,getTime,getTimeForPosts};