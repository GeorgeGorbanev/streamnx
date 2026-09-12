const ho="2634.3.0-external",$c="sample_developer_token",sme="sample_developer_token";
async function Sme(){const n={url:"https://musicstatus.music.apple.com/v1/web/status",headers:{},method:"GET"};return n.headers.Authorization=`Bearer ${$c}`,fetch(n.url,{credentials:"include",...n}).then(i=>i.json()).catch(i=>{Mn("subscriptions").error("error while getting userInfo:",i)})}
