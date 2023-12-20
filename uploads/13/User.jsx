import { useNavigate, useLocation } from "react-router-dom"
import { useEffect } from "react"
import useAxiosPrivate from "../hooks/useAxiosPrivate"
import axios from 'axios'
const User = () => {
    const axiosPrivate = useAxiosPrivate()
    const navigate = useNavigate()
    const location = useLocation()

    useEffect(() => {
        const controller = new AbortController()

        const getUser = async () => {
            try {
                const response = await axiosPrivate.get('/user', {signal: controller.signal})

            } catch (error) {
                if(axios.isCancel(error)){
                    console.log('request cancelled')
                }
                else {
                    console.error(error)          
                    navigate('/login', {state: {from: location}, replace: true})
                }
            }
        }

        getUser()

        return () => {
            controller.abort()
            console.log('unmounted')
        }
    },[])

    return (
        <div>
            <h1>User</h1>
        </div>
    )
}

export default User