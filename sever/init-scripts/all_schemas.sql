--
-- PostgreSQL database cluster dump
--

\restrict 9DErLS04OPyDYgquN959giZUCN4M3Ui6PXOOTj2NaCS5g3bxS5rS7ZTkYvjG8Fs

SET default_transaction_read_only = off;

SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;

--
-- Roles
--

CREATE ROLE root;
ALTER ROLE root WITH SUPERUSER INHERIT CREATEROLE CREATEDB LOGIN REPLICATION BYPASSRLS PASSWORD 'SCRAM-SHA-256$4096:2h95ycrOZEI0ohlFZWNjPg==$JOGQitAgxJvXirmfhzXsO1ka3Ei1FXJvZruiMyDKfgI=:0/ibdp0H9BYiKt2qys7B/7yQWm7XicZzfyX5i3lMSfQ=';

--
-- User Configurations
--








\unrestrict 9DErLS04OPyDYgquN959giZUCN4M3Ui6PXOOTj2NaCS5g3bxS5rS7ZTkYvjG8Fs

--
-- Databases
--

--
-- Database "template1" dump
--

\connect template1

--
-- PostgreSQL database dump
--

\restrict 9mZdQJJOlGMlcbjmTaX7qCYnIJaWcuSTaYEOtIBMCEJAetiR3zLl0u923LzfT3t

-- Dumped from database version 15.17
-- Dumped by pg_dump version 15.17

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- PostgreSQL database dump complete
--

\unrestrict 9mZdQJJOlGMlcbjmTaX7qCYnIJaWcuSTaYEOtIBMCEJAetiR3zLl0u923LzfT3t

--
-- Database "iot_business" dump
--

--
-- PostgreSQL database dump
--

\restrict lTspNazxNYplQNnCaQbOOsuLodG3oACUS3BCqdXLo7bAEfRgzayyA6B1T6c2UPV

-- Dumped from database version 15.17
-- Dumped by pg_dump version 15.17

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: iot_business; Type: DATABASE; Schema: -; Owner: root
--

CREATE DATABASE iot_business WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'en_US.utf8';


ALTER DATABASE iot_business OWNER TO root;

\unrestrict lTspNazxNYplQNnCaQbOOsuLodG3oACUS3BCqdXLo7bAEfRgzayyA6B1T6c2UPV
\connect iot_business
\restrict lTspNazxNYplQNnCaQbOOsuLodG3oACUS3BCqdXLo7bAEfRgzayyA6B1T6c2UPV

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- PostgreSQL database dump complete
--

\unrestrict lTspNazxNYplQNnCaQbOOsuLodG3oACUS3BCqdXLo7bAEfRgzayyA6B1T6c2UPV

--
-- Database "postgres" dump
--

\connect postgres

--
-- PostgreSQL database dump
--

\restrict C0eXDNMsUBGBtc7dFYhY9ioL5mdm4h6LSsazC9prcCtXkY2Rr8zgZfKwT9xcf8z

-- Dumped from database version 15.17
-- Dumped by pg_dump version 15.17

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: iotplus; Type: SCHEMA; Schema: -; Owner: root
--

CREATE SCHEMA iotplus;


ALTER SCHEMA iotplus OWNER TO root;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: access_logs; Type: TABLE; Schema: iotplus; Owner: root
--

CREATE TABLE iotplus.access_logs (
    log_id bigint NOT NULL,
    user_id character varying(50),
    device_id character varying(50) NOT NULL,
    access_time timestamp with time zone DEFAULT now(),
    auth_method character varying(20) NOT NULL,
    result character varying(10) NOT NULL,
    photo_url text,
    reason text,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE iotplus.access_logs OWNER TO root;

--
-- Name: access_logs_log_id_seq; Type: SEQUENCE; Schema: iotplus; Owner: root
--

CREATE SEQUENCE iotplus.access_logs_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE iotplus.access_logs_log_id_seq OWNER TO root;

--
-- Name: access_logs_log_id_seq; Type: SEQUENCE OWNED BY; Schema: iotplus; Owner: root
--

ALTER SEQUENCE iotplus.access_logs_log_id_seq OWNED BY iotplus.access_logs.log_id;


--
-- Name: devices; Type: TABLE; Schema: iotplus; Owner: root
--

CREATE TABLE iotplus.devices (
    device_id character varying(50) NOT NULL,
    location character varying(150) NOT NULL,
    ip_address inet,
    status character varying(20) DEFAULT 'offline'::character varying,
    firmware_version character varying(20),
    last_heartbeat timestamp with time zone,
    created_at timestamp with time zone DEFAULT now()
);


ALTER TABLE iotplus.devices OWNER TO root;

--
-- Name: user_permissions; Type: TABLE; Schema: iotplus; Owner: root
--

CREATE TABLE iotplus.user_permissions (
    user_id character varying(50) NOT NULL,
    name character varying(100) NOT NULL,
    face_feature bytea,
    fingerprint_feature bytea,
    allowed_devices jsonb DEFAULT '[]'::jsonb,
    valid_start timestamp with time zone NOT NULL,
    valid_end timestamp with time zone NOT NULL,
    is_active boolean DEFAULT true,
    updated_at timestamp with time zone DEFAULT now()
);


ALTER TABLE iotplus.user_permissions OWNER TO root;

--
-- Name: admin; Type: TABLE; Schema: public; Owner: root
--

CREATE TABLE public.admin (
    id integer NOT NULL,
    username character varying(50) NOT NULL,
    password character varying(255) NOT NULL,
    status smallint DEFAULT 1 NOT NULL,
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


ALTER TABLE public.admin_id_seq1 OWNER TO root;

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


ALTER TABLE public.message_id_seq OWNER TO root;

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


ALTER TABLE public.node_id_seq OWNER TO root;

--
-- Name: node_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: root
--

ALTER SEQUENCE public.node_id_seq OWNED BY public.node.id;


--
-- Name: access_logs log_id; Type: DEFAULT; Schema: iotplus; Owner: root
--

ALTER TABLE ONLY iotplus.access_logs ALTER COLUMN log_id SET DEFAULT nextval('iotplus.access_logs_log_id_seq'::regclass);


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
-- Name: access_logs access_logs_pkey; Type: CONSTRAINT; Schema: iotplus; Owner: root
--

ALTER TABLE ONLY iotplus.access_logs
    ADD CONSTRAINT access_logs_pkey PRIMARY KEY (log_id);


--
-- Name: devices devices_pkey; Type: CONSTRAINT; Schema: iotplus; Owner: root
--

ALTER TABLE ONLY iotplus.devices
    ADD CONSTRAINT devices_pkey PRIMARY KEY (device_id);


--
-- Name: user_permissions user_permissions_pkey; Type: CONSTRAINT; Schema: iotplus; Owner: root
--

ALTER TABLE ONLY iotplus.user_permissions
    ADD CONSTRAINT user_permissions_pkey PRIMARY KEY (user_id);


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
-- Name: idx_devices_location; Type: INDEX; Schema: iotplus; Owner: root
--

CREATE INDEX idx_devices_location ON iotplus.devices USING btree (location);


--
-- Name: idx_logs_device_time; Type: INDEX; Schema: iotplus; Owner: root
--

CREATE INDEX idx_logs_device_time ON iotplus.access_logs USING btree (device_id, access_time DESC);


--
-- Name: idx_logs_time; Type: INDEX; Schema: iotplus; Owner: root
--

CREATE INDEX idx_logs_time ON iotplus.access_logs USING btree (access_time DESC);


--
-- Name: idx_logs_user_time; Type: INDEX; Schema: iotplus; Owner: root
--

CREATE INDEX idx_logs_user_time ON iotplus.access_logs USING btree (user_id, access_time DESC);


--
-- Name: idx_user_allowed_devices; Type: INDEX; Schema: iotplus; Owner: root
--

CREATE INDEX idx_user_allowed_devices ON iotplus.user_permissions USING gin (allowed_devices);


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

\unrestrict C0eXDNMsUBGBtc7dFYhY9ioL5mdm4h6LSsazC9prcCtXkY2Rr8zgZfKwT9xcf8z

--
-- PostgreSQL database cluster dump complete
--

