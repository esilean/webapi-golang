import React, { Fragment, useEffect, useState } from 'react'

import Input from './form-components/Input'
import TextArea from './form-components/TextArea'
import Select from './form-components/Select'

import Alert from './ui-components/Alert'

import './EditMovie.css'

const mpaaOptions = [
    { id: "G", value: "G" },
    { id: "PG", value: "PG" },
    { id: "PG13", value: "PG13" },
    { id: "R", value: "R" },
    { id: "NC17", value: "NC17" },
]

const EditMovie = (props) => {
    const [movie, setMovie] = useState(
        {
            id: 0,
            title: '',
            release_date: '',
            runtime: '',
            mpaa_rating: '',
            rating: '',
            description: ''
        })
    const [isLoaded, setLoaded] = useState(false)
    const [error, setError] = useState('')
    const [formErrors, setFormErrors] = useState([])
    const [alert, setAlert] = useState({ type: 'd-none', message: '' })

    useEffect(() => {

        const id = props.match.params.id;
        if (id > 0) {
            fetch("http://localhost:4000/v1/movie/" + id)
                .then((response) => {
                    if (response.status !== 200) {
                        setError("Invalid response code: " + response.status)
                        return;
                    }

                    return response.json();
                })
                .then((json) => {

                    const releaseDate = new Date(json.movie.release_date);
                    setMovie({
                        id: id,
                        title: json.movie.title,
                        release_date: releaseDate.toISOString().split("T")[0],
                        runtime: json.movie.runtime,
                        mpaa_rating: json.movie.mpaa_rating,
                        rating: json.movie.rating,
                        description: json.movie.description,
                    })
                })
                .catch((err) => {
                    setAlert({
                        type: 'alert-danger',
                        message: err.message
                    })
                });
        }

        setLoaded(true)

    }, [props.match.params.id])

    function handleSubmit(evt) {
        evt.preventDefault()

        let errors = []
        if (movie.title === '') {
            errors.push('title')
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

        fetch('http://localhost:4000/v1/admin/editmovie', requestOptions)
            .then(response => {
                if (response.status !== 200) {
                    setError("Invalid response code: " + response.status)

                    response.json().then((data) => {

                        setAlert({
                            type: 'alert-danger',
                            message: data.error.message
                        })
                    })

                    return
                }

                setAlert({
                    type: 'alert-success',
                    message: 'Changes saved!'
                })
            })
            .catch((err) => {
                setAlert({
                    type: 'alert-danger',
                    message: err.message
                })
            })
    }

    function handleChange(evt) {
        let value = evt.target.value
        let name = evt.target.name

        setMovie((prev) => ({ ...prev, [name]: value }))
    }

    function hasError(key) {
        return formErrors.indexOf(key) !== -1
    }

    if (!isLoaded) return (<p>Loading...</p>)

    return (
        <Fragment>
            <h2>Add/Edit Movie</h2>
            <Alert
                alertType={alert.type}
                alertMessage={alert.message}
            />
            <hr />
            <form onSubmit={handleSubmit}>
                <input
                    type="hidden"
                    name="id"
                    id="id"
                    value={movie.id}
                    onChange={handleChange}
                />

                <Input
                    title={"Title"}
                    className={hasError('title') ? 'is-invalid' : ''}
                    type={"text"}
                    name={"title"}
                    value={movie.title}
                    handleChange={handleChange}
                    errorDiv={hasError('title') ? 'text-danger' : 'd-none'}
                    errorMsg={'Please enter a title'}
                />

                <Input
                    title={"Release Date"}
                    type={"date"}
                    name={"release_date"}
                    value={movie.release_date}
                    handleChange={handleChange}
                />

                <Input
                    title={"Runtime"}
                    type={"text"}
                    name={"runtime"}
                    value={movie.runtime}
                    handleChange={handleChange}
                />

                <Select
                    title={"MPAA Rating"}
                    name={"mpaa_rating"}
                    options={mpaaOptions}
                    value={movie.mpaa_rating}
                    handleChange={handleChange}
                    placeholder={"Choose..."}
                />

                <Input
                    title={"Rating"}
                    type={"text"}
                    name={"rating"}
                    value={movie.rating}
                    handleChange={handleChange}
                />

                <TextArea
                    title={"Description"}
                    name={"description"}
                    value={movie.description}
                    rows={"3"}
                    handleChange={handleChange}
                />

                <hr />

                <button className='btn btn-primary'>Save</button>
            </form>
        </Fragment>
    );
}

export default EditMovie