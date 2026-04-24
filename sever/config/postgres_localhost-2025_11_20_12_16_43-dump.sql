--
-- PostgreSQL database dump
--

\restrict kp8jWrdclsa4gan1mN2VwwXFr5V7ZJ9kXqNE21gXDPQ1xWVrrnPFW6thr4ScDAO

-- Dumped from database version 18.1 (Debian 18.1-1.pgdg13+2)
-- Dumped by pg_dump version 18.1 (Homebrew)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: admin; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.admin (
    id integer CONSTRAINT admin_id_not_null1 NOT NULL,
    username character varying(50) CONSTRAINT admin_username_not_null1 NOT NULL,
    password character varying(255) CONSTRAINT admin_password_not_null1 NOT NULL,
    status smallint DEFAULT 1 CONSTRAINT admin_status_not_null1 NOT NULL,
    token character varying(255),
    token_expires_at timestamp without time zone,
    token_created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    created_time timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    update_time timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    is_deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.admin OWNER TO root;

--
-- Name: admin_id_seq1; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.admin_id_seq1
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.admin_id_seq1 OWNER TO root;

--
-- Name: admin_id_seq1; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.admin_id_seq1 OWNED BY public.admin.id;


--
-- Name: message; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.message (
    id bigint NOT NULL,
    node_id integer,
    type integer,
    received_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    user_id integer,
    message character(36)
);


ALTER TABLE public.message OWNER TO root;

--
-- Name: message_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.message_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.message_id_seq OWNER TO root;

--
-- Name: message_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.message_id_seq OWNED BY public.message.id;


--
-- Name: node; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.node (
    id integer NOT NULL,
    user_id integer
);


ALTER TABLE public.node OWNER TO root;

--
-- Name: node_id_seq; Type: SEQUENCE; Schema: public; Owner: root
--

CREATE SEQUENCE public.node_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.node_id_seq OWNER TO root;

--
-- Name: node_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.node_id_seq OWNED BY public.node.id;


--
-- Name: admin id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.admin ALTER COLUMN id SET DEFAULT nextval('public.admin_id_seq1'::regclass);


--
-- Name: message id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.message ALTER COLUMN id SET DEFAULT nextval('public.message_id_seq'::regclass);


--
-- Name: node id; Type: DEFAULT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.node ALTER COLUMN id SET DEFAULT nextval('public.node_id_seq'::regclass);


--
-- Data for Name: admin; Type: TABLE DATA; Schema: public; Owner: root
--

INSERT INTO public.admin
(id, username, password, status, token, token_expires_at, token_created_at, created_time, update_time, is_deleted)
VALUES
    (7, 'root', 'admin', 1,
     'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjM2NDkxNDcsInVzZXJuYW1lIjoicm9vdCJ9.WmPb-dKyrsWBBEDRAYg1m6V2b5cs6pEq-tJFA8KjK30',
     '2025-11-20 22:32:27.459736', '2025-11-18 12:31:35.181198',
     '2025-11-18 12:31:35.179963', '2025-11-18 12:31:44.829786', false),

    (5, '赫连国辉', 'aLY6WX_3VlTiDCl', 1,
     NULL, NULL, '2025-11-18 11:07:46.575156',
     '2025-11-18 11:07:46.574745', '2025-11-18 11:07:46.575156', false),

    (2, '羊子豪', 'c1Bzhd4XYvEsgSf', 1,
     NULL, NULL, '2025-11-18 00:27:22.88082',
     '2025-11-18 00:27:22.880081', '2025-11-18 00:27:22.88082', false),

    (6, '许若汐', 'iRFhjDfPd9SIeeZ', 1,
     NULL, NULL, '2025-11-18 11:09:03.935339',
     '2025-11-18 11:09:03.934184', '2025-11-18 11:09:03.935339', false),

    (3, '闾榕融', 'WvNQ0FwONvMlLOf', 1,
     NULL, NULL, '2025-11-18 10:56:49.615499',
     '2025-11-18 10:56:49.615274', '2025-11-18 10:56:49.615499', false),

    (1, 'admin', 'admin', 1,
     'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NjM2OTcyOTMsInVzZXJuYW1lIjoiYWRtaW4ifQ._y5NcTiQ-3vYZdiAfJmUD81wisgRBJy82lpVw6NKpjI',
     '2025-11-21 11:54:53.687216', '2025-11-17 16:24:25.596921',
     '2025-11-17 16:24:25.596921', '2025-11-18 00:26:40.457763', false);


--
-- Data for Name: message; Type: TABLE DATA; Schema: public; Owner: root
--

INSERT INTO public.message (id, node_id, type, received_at, user_id, message) VALUES
                                                                                  (1, 0, 3, '2025-11-18 03:07:08.404684', 1, '23.6'),
                                                                                  (2, 0, 3, '2025-11-18 03:22:33.045693', 1, '15.3'),
                                                                                  (3, 1, 5, '2025-11-18 17:03:12.67259',  1, '8.39'),
                                                                                  (4, 1, 5, '2025-11-18 17:04:03.657486', 1, '-29.95'),
                                                                                  (5, 1, 3, '2025-11-18 17:04:47.692341', 1, '-23.72');


--
-- Data for Name: node; Type: TABLE DATA; Schema: public; Owner: root
--

INSERT INTO public.node (id, user_id) VALUES
                                          (0, 1),
                                          (1, 1),
                                          (2, 1),
                                          (13, 7),
                                          (12, 1);


--
-- Name: admin_id_seq1; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.admin_id_seq1', 7, true);


--
-- Name: message_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.message_id_seq', 5, true);


--
-- Name: node_id_seq; Type: SEQUENCE SET; Schema: public; Owner: root
--

SELECT pg_catalog.setval('public.node_id_seq', 2, true);


--
-- Name: admin admin_pkey; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.admin
    ADD CONSTRAINT admin_pkey PRIMARY KEY (id);


--
-- Name: admin admin_username_key; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.admin
    ADD CONSTRAINT admin_username_key UNIQUE (username);


--
-- Name: message message_pk; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT message_pk PRIMARY KEY (id);


--
-- Name: node node_pk; Type: CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.node
    ADD CONSTRAINT node_pk PRIMARY KEY (id);


--
-- Name: node node_admin_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.node
    ADD CONSTRAINT node_admin_id_fk FOREIGN KEY (user_id) REFERENCES public.admin(id);


--
-- Name: message user_id; Type: FK CONSTRAINT; Schema: public; Owner: root
--

ALTER TABLE ONLY public.message
    ADD CONSTRAINT user_id FOREIGN KEY (user_id) REFERENCES public.admin(id);


--
-- PostgreSQL database dump complete
--

\unrestrict kp8jWrdclsa4gan1mN2VwwXFr5V7ZJ9kXqNE21gXDPQ1xWVrrnPFW6thr4ScDAO

