import React, { Fragment, useState } from 'react'
import Input from './form-components/Input'
import Alert from './ui-components/Alert'

const Login = (props) => {

    const [credentials, setCredentials] = useState(
        {
            email: '',
            password: ''
        })
    const [formErrors, setFormErrors] = useState([])
    const [alertMessage, setAlertMessage] = useState({ type: 'd-none', message: '' })

    function handleChange(evt) {
        let value = evt.target.value
        let name = evt.target.name

        setCredentials((prev) => ({ ...prev, [name]: value }))
    }

    function handleSubmit(evt) {
        evt.preventDefault()

        let errors = []
        if (credentials.email === '') {
            errors.push('email')
        }

        if (credentials.password === '') {
            errors.push('password')
        }

        setFormErrors(errors)

        if (errors.length > 0) {
            return false;
        }

        const data = new FormData(evt.target)
        const payload = Object.fromEntries(data.entries())

        const requestOptions = {
            method: 'POST',
            body: JSON.stringify(payload)
        }

        fetch(`${process.env.REACT_APP_API_URL}/v1/signin`, requestOptions)
            .then(response => {
                return response.json()
            })
            .then((json) => {

                if (json.error) {
                    setAlertMessage({
                        type: 'alert-danger',
                        message: json.error.message
                    })
                    return;
                }

                handleJwtChange(json.token)

                window.localStorage.setItem('token', json.token)

                props.history.push({
                    pathname: '/admin'
                })
            })
            .catch((err) => {
                setAlertMessage({
                    type: 'alert-danger',
                    message: err.message
                })
            })
    }

    function handleJwtChange(jwt) {
        props.handleJwtChange(jwt)
    }

    function hasError(key) {
        return formErrors.indexOf(key) !== -1
    }

    return (
        <Fragment>
            <h2>Login</h2>
            <hr />
            <Alert
                alertType={alertMessage.type}
                alertMessage={alertMessage.message}
            />

            <form className="pt-3" onSubmit={handleSubmit}>
                <Input
                    title={"Email"}
                    type={"email"}
                    name={"email"}
                    handleChange={handleChange}
                    className={hasError("email") ? "is-invalid" : ""}
                    errorDiv={hasError("email") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a valid email address"}
                />

                <Input
                    title={"Password"}
                    type={"password"}
                    name={"password"}
                    handleChange={handleChange}
                    className={hasError("password") ? "is-invalid" : ""}
                    errorDiv={hasError("password") ? "text-danger" : "d-none"}
                    errorMsg={"Please enter a password"}
                />

                <hr />
                <button className="btn btn-primary">Login</button>
            </form>
        </Fragment>
    )
}

export default Login