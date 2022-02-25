import React, { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import Input from "./form-components/Input";
import Alert from "./ui-components/Alert";

const MoviesGraphQL = () => {
    const [movies, setMovies] = useState([])
    const [isLoaded, setLoaded] = useState(false)
    const [alertMessage, setAlertMessage] = useState({ type: 'd-none', message: '' })
    const [searchTerm, setSearchTerm] = useState('')

    useEffect(() => {

        const payload = `
            {
                list {
                    id
                    title
                    runtime
                    year,
                    description
                }
            }
        `

        const headers = new Headers()
        headers.append("Content-Type", "application/json")

        const requestOptions = {
            method: 'POST',
            body: payload,
            headers: headers
        }

        fetch(`${process.env.REACT_APP_API_URL}/v1/graphql`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                return Object.values(data.data.list)
            })
            .then((movies) => {
                setMovies(movies)
            })
            .catch((err) => {
                setAlertMessage({
                    type: 'alert-danger',
                    message: err.message
                })
            })
            .finally(() => {
                setLoaded(true)
            })

    }, [])

    useEffect(() => {
        if (searchTerm.length === 0 || searchTerm.length >= 3) {
            performSearch(searchTerm)
        } else {
            setMovies([])
        }
    }, [searchTerm])

    function handleChange(evt) {
        setSearchTerm(evt.target.value)
    }


    function performSearch(searchTermInput) {

        const payload = `
            {
                search(titleContains: "${searchTermInput}") {
                    id
                    title
                    runtime
                    year,
                    description
                }
            }
        `

        const headers = new Headers()
        headers.append("Content-Type", "application/json")

        const requestOptions = {
            method: 'POST',
            body: payload,
            headers: headers
        }

        fetch(`${process.env.REACT_APP_API_URL}/v1/graphql`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                return Object.values(data.data.search)
            })
            .then((movies) => {
                if (movies.length > 0) {
                    setMovies(movies)
                } else {
                    setMovies([])
                }
            })
            .catch((err) => {
                setAlertMessage({
                    type: 'alert-danger',
                    message: err.message
                })
            })
            .finally(() => {
                setLoaded(true)
            })
    }

    if (!isLoaded) return (<p>Loading...</p>)

    return (
        <>
            <h2>Choose a movie</h2>
            <Alert
                alertType={alertMessage.type}
                alertMessage={alertMessage.message}
            />
            <hr />

            <Input
                title={'Search'}
                type={'text'}
                name={'search'}
                value={searchTerm}
                handleChange={handleChange}
            />

            <div className='list-group'>
                {movies.map((m) => (
                    <Link key={m.id} to={`/moviesgraphql/${m.id}`} className="list-group-item list-group-item-action">
                        <strong>{m.title}</strong> <br />
                        <small>
                            ({m.year}) - {m.runtime} minutes
                        </small>
                        <br />
                        {m.description.slice(0, 100)}...
                    </Link>
                ))}
            </div>
        </>
    )

}

export default MoviesGraphQL