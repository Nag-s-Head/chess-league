--
-- PostgreSQL database dump
--


-- Dumped from database version 18.2 (Debian 18.2-1.pgdg13+1)
-- Dumped by pg_dump version 18.2 (Debian 18.2-1.pgdg13+1)




--
-- Name: admin_users; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.admin_users (
    id uuid NOT NULL,
    name text NOT NULL,
    oauth_id text NOT NULL,
    created timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    session_key text,
    last_login timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_ip text,
    last_user_agent text
);


--
-- Name: audit_log_game_affected; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.audit_log_game_affected (
    audit_log_id uuid NOT NULL,
    game_ikey bigint NOT NULL
);


--
-- Name: audit_log_player_affected; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.audit_log_player_affected (
    audit_log_id uuid NOT NULL,
    player_id text NOT NULL,
    elo_change integer NOT NULL
);


--
-- Name: audit_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.audit_logs (
    id uuid NOT NULL,
    created timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    done_by uuid NOT NULL,
    operation_name text NOT NULL,
    operation_description text NOT NULL
);


--
-- Name: game_ikey_sequence; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.game_ikey_sequence
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: games; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.games (
    player_white text NOT NULL,
    player_black text NOT NULL,
    score text NOT NULL,
    submitter text NOT NULL,
    played timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    elo_given integer NOT NULL,
    elo_taken integer NOT NULL,
    submit_ip text,
    submit_user_agent text,
    ikey bigint NOT NULL,
    CONSTRAINT games_score_check CHECK (((score = '1-0'::text) OR (score = '0-1'::text) OR (score = '1/2-1/2'::text)))
);


--
-- Name: migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.migrations (
    version integer NOT NULL,
    date timestamp with time zone NOT NULL
);


--
-- Name: players; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.players (
    id text NOT NULL,
    name text NOT NULL,
    name_normalised text NOT NULL,
    elo integer DEFAULT 1000,
    join_time timestamp with time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    CONSTRAINT players_elo_check CHECK ((elo >= 0))
);


--
-- Data for Name: admin_users; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.admin_users VALUES ('26b9950d-6473-476a-b833-26c4f4098ae9', 'Danny Piper', 'djpiper28', '2026-03-18 18:19:00.714375+00', 'LGSN4UH56X3T5CWY7FUC2WWLV5LVKUXSQOPPQ7UJJFZVFOX25HBGOGRA55EYXAVP3LQ7BJG7S7PT5NHEW5Y74RVSHBX5EMSL62YDUGU3WFEQAOLEHKWQEGICOJRV6KBLXH', '2026-03-25 14:27:08.544765+00', '10.244.0.1:42027', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0');
INSERT INTO public.admin_users VALUES ('3295095f-2dc9-4881-b5e3-eb6abea2b62c', 'Rhys John', 'rhysajohn', '2026-03-19 16:29:58.470574+00', 'PKETWJCGTEDDUW6LAHXE2DUGQJABJTLTZ4JIYSPSM4UTWRT5QFB6HQV3TFGRUYYM5RIGXMEVZKRIYTB3QK2SAHR5NAD3HZIPIFMVWZDCQQNGPSMXP2NL5UQU54HJ47PPBP', '2026-03-31 23:43:35.33587+00', '10.244.0.1:48538', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1');


--
-- Data for Name: audit_log_game_affected; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: audit_log_player_affected; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: audit_logs; Type: TABLE DATA; Schema: public; Owner: -
--



--
-- Data for Name: games; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-02-27 19:46:10.294196+00', false, 15, -15, '10.244.0.1:49333', 'Mozilla/5.0 (Android 16; Mobile; rv:147.0) Gecko/147.0 Firefox/147.0', 1);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'ca8d7325-406e-4742-afa3-f8404f155139', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-02-27 19:47:35.085186+00', false, 16, -16, '10.244.0.1:36001', 'Mozilla/5.0 (Android 16; Mobile; rv:147.0) Gecko/147.0 Firefox/147.0', 2);
INSERT INTO public.games VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', 'ca8d7325-406e-4742-afa3-f8404f155139', '2026-02-27 20:17:01.576788+00', false, 14, -14, '10.244.0.1:19589', 'Mozilla/5.0 (Android 16; Mobile; rv:147.0) Gecko/147.0 Firefox/147.0', 3);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-02-27 20:23:43.694993+00', false, 15, -15, '10.244.0.1:12691', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/145.0.7632.108 Mobile/15E148 Safari/604.1', 4);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', 'ca8d7325-406e-4742-afa3-f8404f155139', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-02-27 20:26:26.872319+00', false, 18, -18, '10.244.0.1:32131', 'Mozilla/5.0 (Android 16; Mobile; rv:147.0) Gecko/147.0 Firefox/147.0', 5);
INSERT INTO public.games VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '2026-02-27 20:45:04.162721+00', false, 15, -15, '10.244.0.1:8757', 'Mozilla/5.0 (Android 16; Mobile; rv:147.0) Gecko/147.0 Firefox/147.0', 6);
INSERT INTO public.games VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', 'ca8d7325-406e-4742-afa3-f8404f155139', '2026-02-27 20:57:31.67992+00', false, 15, -15, '10.244.0.1:58681', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.3 Mobile/15E148 Safari/604.1', 10);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-02-27 21:01:58.832701+00', false, 14, -14, '10.244.0.1:31311', 'Mozilla/5.0 (Android 16; Mobile; rv:147.0) Gecko/147.0 Firefox/147.0', 11);
INSERT INTO public.games VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '2026-02-27 21:22:24.896535+00', false, 12, -12, '10.244.0.1:53419', 'Mozilla/5.0 (Android 16; Mobile; rv:147.0) Gecko/147.0 Firefox/147.0', 16);
INSERT INTO public.games VALUES ('1c401b80-b664-4f1c-9fd9-6e97c1282135', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '1c401b80-b664-4f1c-9fd9-6e97c1282135', '2026-02-27 21:27:12.238712+00', false, 15, -15, '10.244.0.1:63917', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/145.0.7632.108 Mobile/15E148 Safari/604.1', 9);
INSERT INTO public.games VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '2026-02-27 21:40:57.971918+00', false, 17, -17, '10.244.0.1:64697', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/145.0.7632.108 Mobile/15E148 Safari/604.1', 18);
INSERT INTO public.games VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '2026-02-27 21:41:03.573897+00', false, 15, -15, '10.244.0.1:5841', 'Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/29.0 Chrome/136.0.0.0 Mobile Safari/537.36', 17);
INSERT INTO public.games VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', 'ca8d7325-406e-4742-afa3-f8404f155139', '2026-02-27 22:05:42.490584+00', false, 15, -15, '10.244.0.1:54678', 'Mozilla/5.0 (Android 16; Mobile; rv:147.0) Gecko/147.0 Firefox/147.0', 19);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-06 19:56:12.993293+00', false, 18, -18, '10.244.0.1:1864', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 21);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-06 20:04:08.557554+00', false, 14, -14, '10.244.0.1:8252', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 22);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-06 20:16:23.611019+00', false, 17, -17, '10.244.0.1:24625', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 23);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-06 21:15:33.295811+00', false, 14, -14, '10.244.0.1:21086', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 24);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-06 21:32:16.02213+00', false, 17, -17, '10.244.0.1:31423', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 25);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-06 21:44:59.668417+00', false, 14, -14, '10.244.0.1:27797', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 26);
INSERT INTO public.games VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', 'ca8d7325-406e-4742-afa3-f8404f155139', '2026-03-06 22:08:34.365292+00', false, 13, -13, '10.244.0.1:4955', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 27);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-06 22:23:41.463536+00', false, 14, -14, '10.244.0.1:16925', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 28);
INSERT INTO public.games VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', 'ca8d7325-406e-4742-afa3-f8404f155139', '2026-03-06 22:54:21.354472+00', false, 11, -11, '10.244.0.1:55179', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 29);
INSERT INTO public.games VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '2026-03-06 23:27:34.23544+00', false, 14, -14, '10.244.0.1:44924', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 30);
INSERT INTO public.games VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '2026-03-06 23:50:24.625416+00', false, 16, -16, '10.244.0.1:24929', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 31);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-07 00:00:09.120895+00', false, 12, -12, '10.244.0.1:16542', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 32);
INSERT INTO public.games VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '2026-03-07 00:02:39.331579+00', false, 16, -16, '10.244.0.1:59942', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/145.0.7632.108 Mobile/15E148 Safari/604.1', 33);
INSERT INTO public.games VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '2026-03-07 00:03:15.629716+00', false, 14, -14, '10.244.0.1:24824', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/145.0.7632.108 Mobile/15E148 Safari/604.1', 34);
INSERT INTO public.games VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '0-1', 'ca8d7325-406e-4742-afa3-f8404f155139', '2026-03-07 00:27:07.847761+00', false, 11, -11, '10.244.0.1:20452', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.3 Mobile/15E148 Safari/604.1', 35);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-07 00:31:31.025723+00', false, 16, -16, '10.244.0.1:9746', 'Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) SamsungBrowser/29.0 Chrome/136.0.0.0 Mobile Safari/537.36', 36);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 18:10:35.584169+00', false, 19, -19, '10.244.0.1:7914', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 38);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 18:16:54.174035+00', false, 17, -17, '10.244.0.1:49201', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 39);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 18:29:22.243488+00', false, 16, -16, '10.244.0.1:29577', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 40);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 18:42:00.914532+00', false, 15, -15, '10.244.0.1:46870', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 41);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 18:52:24.545361+00', false, 14, -14, '10.244.0.1:57606', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 42);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 19:16:19.993549+00', false, 17, -17, '10.244.0.1:29357', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 43);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 19:24:48.266473+00', false, 14, -14, '10.244.0.1:1663', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 44);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 19:50:27.153492+00', false, 17, -17, '10.244.0.1:29566', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 45);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 20:04:54.324373+00', false, 15, -15, '10.244.0.1:27124', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 46);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 20:16:33.248527+00', false, 14, -14, '10.244.0.1:24519', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 47);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 20:16:46.048597+00', false, 17, -17, '10.244.0.1:22753', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 48);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 20:16:48.364507+00', false, 16, -16, '10.244.0.1:10931', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 49);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 20:38:32.940764+00', false, 14, -14, '10.244.0.1:21851', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 50);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 20:45:54.221161+00', false, 17, -17, '10.244.0.1:43502', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 51);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 20:50:37.215622+00', false, 15, -15, '10.244.0.1:15079', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 52);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 20:59:43.592072+00', false, 17, -17, '10.244.0.1:17875', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 53);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 21:12:09.140716+00', false, 15, -15, '10.244.0.1:10572', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 54);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 21:27:04.176343+00', false, 16, -16, '10.244.0.1:56434', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 55);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 21:43:49.612604+00', false, 15, -15, '10.244.0.1:19914', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 56);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1/2-1/2', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 22:25:35.607931+00', false, 1, -1, '10.244.0.1:59393', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 57);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 22:38:00.546852+00', false, 16, -16, '10.244.0.1:51509', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 58);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 22:44:49.866552+00', false, 15, -15, '10.244.0.1:2057', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 59);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 23:18:46.70365+00', false, 16, -16, '10.244.0.1:30451', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 60);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-09 23:32:06.430692+00', false, 16, -16, '10.244.0.1:6259', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 61);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-09 23:50:38.065346+00', false, 14, -14, '10.244.0.1:26531', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 62);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'ca8d7325-406e-4742-afa3-f8404f155139', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-12 20:02:06.993542+00', false, 12, -12, '10.244.0.1:9309', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/145.0.7632.108 Mobile/15E148 Safari/604.1', 63);
INSERT INTO public.games VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', 'ca8d7325-406e-4742-afa3-f8404f155139', '2026-03-12 20:07:27.514713+00', false, 11, -11, '10.244.0.1:5579', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/145.0.7632.108 Mobile/15E148 Safari/604.1', 64);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-12 20:24:08.616037+00', false, 16, -16, '10.244.0.1:21242', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 66);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-12 20:36:36.33178+00', false, 15, -15, '10.244.0.1:26661', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 67);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-12 21:09:06.808845+00', false, 16, -16, '10.244.0.1:42844', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 68);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-12 21:50:20.907952+00', false, 15, -15, '10.244.0.1:23259', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 69);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1/2-1/2', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-12 22:13:07.445793+00', false, 1, -1, '10.244.0.1:63535', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 71);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', 'ca8d7325-406e-4742-afa3-f8404f155139', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-12 22:32:55.403091+00', false, 9, -9, '10.244.0.1:1972', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 72);
INSERT INTO public.games VALUES ('c2d416c4-f457-4b62-a9d4-cdf7c2c64624', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', 'c2d416c4-f457-4b62-a9d4-cdf7c2c64624', '2026-03-14 17:43:48.709703+00', false, 17, -17, '10.244.0.1:53952', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 74);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', 'c2d416c4-f457-4b62-a9d4-cdf7c2c64624', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-14 17:52:22.74699+00', false, 15, -15, '10.244.0.1:16327', 'Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/145.0.0.0 Mobile Safari/537.36', 75);
INSERT INTO public.games VALUES ('c2d416c4-f457-4b62-a9d4-cdf7c2c64624', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', 'c2d416c4-f457-4b62-a9d4-cdf7c2c64624', '2026-03-14 18:03:45.940161+00', false, 17, -17, '10.244.0.1:59778', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 76);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-14 23:54:19.383954+00', false, 15, -15, '10.244.0.1:60251', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 77);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-14 23:54:38.205628+00', false, 17, -17, '10.244.0.1:57325', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 78);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-14 23:54:41.076381+00', false, 15, -15, '10.244.0.1:36459', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 79);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'ca8d7325-406e-4742-afa3-f8404f155139', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-14 23:54:48.31723+00', false, 9, -9, '10.244.0.1:44641', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 80);
INSERT INTO public.games VALUES ('d19e9c33-643d-488e-8c35-c01974bcafff', 'a3a41684-a197-4aad-9521-c9ded1fc2ebe', '1-0', 'd19e9c33-643d-488e-8c35-c01974bcafff', '2026-03-16 20:03:29.453882+00', false, 15, -15, '10.244.0.1:27694', 'Mozilla/5.0 (iPhone; CPU iPhone OS 26_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 81);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-16 20:04:51.673465+00', false, 13, -13, '10.244.0.1:12424', 'Mozilla/5.0 (iPhone; CPU iPhone OS 26_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 82);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', 'a3a41684-a197-4aad-9521-c9ded1fc2ebe', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-16 20:15:40.273402+00', false, 15, -15, '10.244.0.1:42226', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 83);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-16 21:11:34.472814+00', false, 13, -13, '10.244.0.1:36884', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 84);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', 'd19e9c33-643d-488e-8c35-c01974bcafff', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-16 21:42:32.906758+00', false, 14, -14, '10.244.0.1:10794', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 85);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'd19e9c33-643d-488e-8c35-c01974bcafff', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-16 22:13:10.918208+00', false, 13, -13, '10.244.0.1:61132', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 86);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-16 22:51:25.434322+00', false, 11, -11, '10.244.0.1:5771', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 88);
INSERT INTO public.games VALUES ('d19e9c33-643d-488e-8c35-c01974bcafff', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', 'd19e9c33-643d-488e-8c35-c01974bcafff', '2026-03-16 23:11:52.83131+00', false, 13, -13, '10.244.0.1:63905', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 90);
INSERT INTO public.games VALUES ('a3a41684-a197-4aad-9521-c9ded1fc2ebe', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', 'a3a41684-a197-4aad-9521-c9ded1fc2ebe', '2026-03-16 23:12:40.649276+00', false, 15, -15, '10.244.0.1:49294', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 91);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', 'a3a41684-a197-4aad-9521-c9ded1fc2ebe', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-16 23:12:47.719131+00', false, 14, -14, '10.244.0.1:59041', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 92);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'a3a41684-a197-4aad-9521-c9ded1fc2ebe', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-16 23:40:25.795424+00', false, 9, -9, '10.244.0.1:21664', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 89);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1/2-1/2', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 20:21:22.958144+00', false, 5, -5, '10.244.0.1:40077', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 94);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-17 20:35:32.020004+00', false, 19, -19, '10.244.0.1:47915', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 95);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'ca8d7325-406e-4742-afa3-f8404f155139', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 21:04:42.901885+00', false, 7, -7, '10.244.0.1:28812', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.3 Mobile/15E148 Safari/604.1', 96);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-17 21:13:39.978772+00', false, 18, -18, '10.244.0.1:25666', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 97);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 21:13:48.76333+00', false, 13, -13, '10.244.0.1:8400', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 98);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', 'ca8d7325-406e-4742-afa3-f8404f155139', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-17 21:13:58.327666+00', false, 9, -9, '10.244.0.1:53597', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 99);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 21:28:04.795844+00', false, 13, -13, '10.244.0.1:45187', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 100);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-17 21:38:26.06853+00', false, 18, -18, '10.244.0.1:10418', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 101);
INSERT INTO public.games VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', 'ca8d7325-406e-4742-afa3-f8404f155139', '2026-03-17 22:03:06.893341+00', false, 7, -7, '10.244.0.1:4643', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.3 Mobile/15E148 Safari/604.1', 102);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 22:12:55.106878+00', false, 17, -17, '10.244.0.1:41004', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 103);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 22:28:41.649044+00', false, 14, -14, '10.244.0.1:15147', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 104);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 23:01:35.084615+00', false, 13, -13, '10.244.0.1:27946', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 105);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 23:02:10.900567+00', false, 18, -18, '10.244.0.1:7004', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 106);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'ca8d7325-406e-4742-afa3-f8404f155139', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 23:02:32.500857+00', false, 7, -7, '10.244.0.1:57440', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 107);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-17 23:22:52.27739+00', false, 13, -13, '10.244.0.1:58226', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 108);
INSERT INTO public.games VALUES ('9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', 'f4836c81-265f-4ea0-9fff-f8a4692daf6a', '1-0', '9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '2026-03-18 17:18:04.339264+00', false, 15, -15, '10.244.0.1:6792', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 109);
INSERT INTO public.games VALUES ('e61d2db3-71e8-4196-963c-2f0f761f5432', 'f4836c81-265f-4ea0-9fff-f8a4692daf6a', '0-1', 'e61d2db3-71e8-4196-963c-2f0f761f5432', '2026-03-18 17:19:33.235949+00', false, 16, -16, '10.244.0.1:36447', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 110);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'f4836c81-265f-4ea0-9fff-f8a4692daf6a', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-18 17:32:56.693814+00', false, 19, -19, '10.244.0.1:26925', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 111);
INSERT INTO public.games VALUES ('f4836c81-265f-4ea0-9fff-f8a4692daf6a', '9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '0-1', 'f4836c81-265f-4ea0-9fff-f8a4692daf6a', '2026-03-18 18:28:34.916641+00', false, 15, -15, '10.244.0.1:58433', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 112);
INSERT INTO public.games VALUES ('9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '2026-03-18 19:05:16.85312+00', false, 17, -17, '10.244.0.1:44579', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 113);
INSERT INTO public.games VALUES ('e61d2db3-71e8-4196-963c-2f0f761f5432', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', 'e61d2db3-71e8-4196-963c-2f0f761f5432', '2026-03-19 15:22:37.021412+00', false, 13, -13, '10.244.0.1:7845', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/26.3 Mobile/15E148 Safari/604.1', 117);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-19 15:47:38.301784+00', false, 14, -14, '10.244.0.1:10256', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 115);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 16:01:05.096067+00', false, 17, -17, '10.244.0.1:64877', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 116);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 16:21:16.661239+00', false, 14, -14, '10.244.0.1:10969', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 118);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 16:37:42.543365+00', false, 13, -13, '10.244.0.1:3179', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 119);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 16:57:25.897596+00', false, 18, -18, '10.244.0.1:10483', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 120);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 17:02:23.901867+00', false, 16, -16, '10.244.0.1:34632', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 121);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-19 17:08:43.801052+00', false, 15, -15, '10.244.0.1:62403', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 122);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 17:21:10.953372+00', false, 16, -16, '10.244.0.1:9142', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 123);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 17:21:19.298579+00', false, 15, -15, '10.244.0.1:1088', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 124);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 17:21:25.83993+00', false, 14, -14, '10.244.0.1:56327', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 125);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-19 17:30:49.254893+00', false, 13, -13, '10.244.0.1:55389', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 126);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-19 18:13:39.617087+00', false, 12, -12, '10.244.0.1:10552', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 127);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 18:22:44.987427+00', false, 19, -19, '10.244.0.1:28752', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 128);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 18:35:00.029264+00', false, 12, -12, '10.244.0.1:22346', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 129);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-19 18:35:45.25983+00', false, 11, -11, '10.244.0.1:5078', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 130);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 18:42:14.791929+00', false, 20, -20, '10.244.0.1:35584', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 131);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-19 19:02:34.528968+00', false, 12, -12, '10.244.0.1:61116', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 132);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-19 19:15:51.953522+00', false, 19, -19, '10.244.0.1:60083', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 133);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-19 19:36:50.45703+00', false, 12, -12, '10.244.0.1:45582', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 134);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '0-1', '9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '2026-03-19 20:03:12.711938+00', false, 17, -17, '10.244.0.1:46599', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 136);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-19 21:06:05.092777+00', false, 15, -15, '10.244.0.1:51693', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.40 Mobile/15E148 Safari/604.1', 137);
INSERT INTO public.games VALUES ('d5a92e7c-9a22-4ba3-9bef-d273665d2c68', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', 'd5a92e7c-9a22-4ba3-9bef-d273665d2c68', '2026-03-19 21:45:33.441063+00', false, 19, -19, '10.244.0.1:55275', 'Mozilla/5.0 (Android 13; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 138);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'd5a92e7c-9a22-4ba3-9bef-d273665d2c68', '1/2-1/2', 'd5a92e7c-9a22-4ba3-9bef-d273665d2c68', '2026-03-19 22:39:46.132502+00', false, 2, -2, '10.244.0.1:33334', 'Mozilla/5.0 (Android 13; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 139);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 17:25:39.14048+00', false, 18, -18, '10.244.0.1:49699', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 140);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 17:41:10.877932+00', false, 14, -14, '10.244.0.1:25518', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 141);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 18:01:20.326141+00', false, 13, -13, '10.244.0.1:35867', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 142);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1/2-1/2', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 18:22:05.728985+00', false, 3, -3, '10.244.0.1:46612', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 143);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 18:33:03.880634+00', false, 12, -12, '10.244.0.1:25066', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 144);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 18:53:00.975626+00', false, 11, -11, '10.244.0.1:40438', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 145);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 19:11:03.182899+00', false, 10, -10, '10.244.0.1:50150', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 146);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 19:17:19.52264+00', false, 9, -9, '10.244.0.1:32097', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 147);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 19:30:54.32176+00', false, 9, -9, '10.244.0.1:58148', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 148);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 19:44:43.614075+00', false, 22, -22, '10.244.0.1:21378', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 149);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 19:54:42.801644+00', false, 10, -10, '10.244.0.1:24273', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 150);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 20:13:18.618573+00', false, 21, -21, '10.244.0.1:49693', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 151);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 20:24:27.234402+00', false, 10, -10, '10.244.0.1:50040', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 152);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 20:40:22.962524+00', false, 20, -20, '10.244.0.1:39318', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 153);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 20:49:17.442903+00', false, 19, -19, '10.244.0.1:1300', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 154);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 20:59:30.795053+00', false, 13, -13, '10.244.0.1:57237', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 155);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 21:22:24.219186+00', false, 12, -12, '10.244.0.1:36325', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 156);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 21:31:17.765892+00', false, 19, -19, '10.244.0.1:2095', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 157);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 21:31:42.680927+00', false, 12, -12, '10.244.0.1:39268', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 158);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 21:31:45.870532+00', false, 11, -11, '10.244.0.1:49419', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 159);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 22:01:22.864639+00', false, 12, -12, '10.244.0.1:1425', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 160);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 22:10:54.20117+00', false, 10, -10, '10.244.0.1:49281', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 161);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 22:23:43.152192+00', false, 21, -21, '10.244.0.1:37788', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 162);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 22:30:38.134611+00', false, 11, -11, '10.244.0.1:5660', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 163);
INSERT INTO public.games VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 22:37:15.163829+00', false, 10, -10, '10.244.0.1:62643', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 164);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '0-1', '57735a42-9df3-495f-a037-ea8fc0a94331', '2026-03-24 22:54:58.281311+00', false, 21, -21, '10.244.0.1:43868', 'Mozilla/5.0 (Android 16; Mobile; rv:148.0) Gecko/148.0 Firefox/148.0', 165);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-29 21:26:39.837694+00', false, 11, -11, '10.244.0.1:49909', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 166);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '57735a42-9df3-495f-a037-ea8fc0a94331', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-29 22:13:48.720042+00', false, 10, -10, '10.244.0.1:3811', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 167);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'd19e9c33-643d-488e-8c35-c01974bcafff', '0-1', 'd19e9c33-643d-488e-8c35-c01974bcafff', '2026-03-31 20:13:53.209179+00', false, 18, -18, '10.244.0.1:28451', 'Mozilla/5.0 (iPhone; CPU iPhone OS 26_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 168);
INSERT INTO public.games VALUES ('d19e9c33-643d-488e-8c35-c01974bcafff', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', 'd19e9c33-643d-488e-8c35-c01974bcafff', '2026-03-31 20:57:48.387854+00', false, 17, -17, '10.244.0.1:18689', 'Mozilla/5.0 (iPhone; CPU iPhone OS 26_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 169);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', '9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', '1-0', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-31 21:03:33.123797+00', false, 15, -15, '10.244.0.1:44517', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 170);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'd19e9c33-643d-488e-8c35-c01974bcafff', '0-1', 'd19e9c33-643d-488e-8c35-c01974bcafff', '2026-03-31 21:27:26.0326+00', false, 16, -16, '10.244.0.1:1097', 'Mozilla/5.0 (iPhone; CPU iPhone OS 26_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 171);
INSERT INTO public.games VALUES ('d19e9c33-643d-488e-8c35-c01974bcafff', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '1-0', 'd19e9c33-643d-488e-8c35-c01974bcafff', '2026-03-31 22:07:54.64079+00', false, 14, -14, '10.244.0.1:17251', 'Mozilla/5.0 (iPhone; CPU iPhone OS 26_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 172);
INSERT INTO public.games VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'd19e9c33-643d-488e-8c35-c01974bcafff', '0-1', 'd19e9c33-643d-488e-8c35-c01974bcafff', '2026-03-31 22:20:59.371507+00', false, 13, -13, '10.244.0.1:37553', 'Mozilla/5.0 (iPhone; CPU iPhone OS 26_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 173);
INSERT INTO public.games VALUES ('d19e9c33-643d-488e-8c35-c01974bcafff', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '0-1', '808100a8-9d64-4c0b-948c-905eb4c7dd5d', '2026-03-31 22:52:19.299385+00', false, 18, -18, '10.244.0.1:28723', 'Mozilla/5.0 (iPhone; CPU iPhone OS 18_7_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/146.0.7680.151 Mobile/15E148 Safari/604.1', 174);


--
-- Data for Name: migrations; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.migrations VALUES (0, '2026-02-24 20:36:32.844199+00');
INSERT INTO public.migrations VALUES (1, '2026-02-25 23:44:36.861414+00');
INSERT INTO public.migrations VALUES (2, '2026-02-27 16:58:36.75004+00');
INSERT INTO public.migrations VALUES (3, '2026-02-27 16:58:36.760357+00');
INSERT INTO public.migrations VALUES (4, '2026-03-03 21:38:43.811626+00');
INSERT INTO public.migrations VALUES (5, '2026-03-03 21:38:43.815439+00');
INSERT INTO public.migrations VALUES (6, '2026-03-03 22:24:03.450771+00');
INSERT INTO public.migrations VALUES (7, '2026-03-18 18:17:22.213581+00');
INSERT INTO public.migrations VALUES (8, '2026-03-18 23:47:26.055367+00');
INSERT INTO public.migrations VALUES (9, '2026-03-18 23:47:26.115615+00');


--
-- Data for Name: players; Type: TABLE DATA; Schema: public; Owner: -
--

INSERT INTO public.players VALUES ('d19e9c33-643d-488e-8c35-c01974bcafff', 'Rob Hornby', 'rob hornby', 1089, '2026-03-16 20:03:29.449525+00', false);
INSERT INTO public.players VALUES ('808100a8-9d64-4c0b-948c-905eb4c7dd5d', 'Rhys', 'rhys', 1056, '2026-02-27 19:46:10.292818+00', false);
INSERT INTO public.players VALUES ('a3a41684-a197-4aad-9521-c9ded1fc2ebe', 'Tom Simmonds', 'tom simmonds', 932, '2026-03-16 20:03:29.453439+00', false);
INSERT INTO public.players VALUES ('64a4a94a-951f-4d1c-8aa7-7fbe4268a2c4', 'Dan', 'dan', 1063, '2026-02-27 20:23:43.694445+00', false);
INSERT INTO public.players VALUES ('1c401b80-b664-4f1c-9fd9-6e97c1282135', 'Matt', 'matt', 985, '2026-02-27 21:27:12.237943+00', false);
INSERT INTO public.players VALUES ('c2d416c4-f457-4b62-a9d4-cdf7c2c64624', 'Callum', 'callum', 1019, '2026-03-14 17:43:48.708766+00', false);
INSERT INTO public.players VALUES ('ca8d7325-406e-4742-afa3-f8404f155139', 'Adam', 'adam', 852, '2026-02-27 19:47:35.08484+00', false);
INSERT INTO public.players VALUES ('d5a92e7c-9a22-4ba3-9bef-d273665d2c68', 'Rory', 'rory', 1021, '2026-03-19 21:45:33.440011+00', false);
INSERT INTO public.players VALUES ('57735a42-9df3-495f-a037-ea8fc0a94331', 'Danny Piper', 'danny piper', 961, '2026-02-27 19:46:10.293881+00', false);
INSERT INTO public.players VALUES ('f4836c81-265f-4ea0-9fff-f8a4692daf6a', 'Dave Evans', 'dave evans', 1005, '2026-03-18 17:18:04.33882+00', false);
INSERT INTO public.players VALUES ('e61d2db3-71e8-4196-963c-2f0f761f5432', 'Spencer Marlow', 'spencer marlow', 971, '2026-03-18 17:19:33.234699+00', false);
INSERT INTO public.players VALUES ('9905b7ff-edd0-4c8d-b22e-9ff5cab1f74f', 'Martyn', 'martyn', 1046, '2026-03-18 17:18:04.337734+00', false);


--
-- Name: game_ikey_sequence; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.game_ikey_sequence', 175, true);


--
-- Name: admin_users admin_users_oauth_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_oauth_id_key UNIQUE (oauth_id);


--
-- Name: admin_users admin_users_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.admin_users
    ADD CONSTRAINT admin_users_pkey PRIMARY KEY (id);


--
-- Name: audit_log_game_affected audit_log_game_affected_audit_log_id_game_ikey_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_log_game_affected
    ADD CONSTRAINT audit_log_game_affected_audit_log_id_game_ikey_key UNIQUE (audit_log_id, game_ikey);


--
-- Name: audit_log_player_affected audit_log_player_affected_audit_log_id_player_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_log_player_affected
    ADD CONSTRAINT audit_log_player_affected_audit_log_id_player_id_key UNIQUE (audit_log_id, player_id);


--
-- Name: audit_logs audit_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_pkey PRIMARY KEY (id);


--
-- Name: games games_ikey_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_ikey_key UNIQUE (ikey);


--
-- Name: migrations migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.migrations
    ADD CONSTRAINT migrations_pkey PRIMARY KEY (version);


--
-- Name: players players_name_normalised_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_name_normalised_key UNIQUE (name_normalised);


--
-- Name: players players_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.players
    ADD CONSTRAINT players_pkey PRIMARY KEY (id);


--
-- Name: idx_admin_users_session_key; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_admin_users_session_key ON public.admin_users USING btree (session_key);


--
-- Name: idx_audit_log_game_affected_audit_log_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_log_game_affected_audit_log_id ON public.audit_log_game_affected USING btree (audit_log_id);


--
-- Name: idx_audit_log_game_affected_game_ikey; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_log_game_affected_game_ikey ON public.audit_log_game_affected USING btree (game_ikey);


--
-- Name: idx_audit_log_player_affected_audit_log_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_log_player_affected_audit_log_id ON public.audit_log_player_affected USING btree (audit_log_id);


--
-- Name: idx_audit_log_player_affected_player_id; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_audit_log_player_affected_player_id ON public.audit_log_player_affected USING btree (player_id);


--
-- Name: idx_games_played; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_games_played ON public.games USING btree (played);


--
-- Name: idx_games_player_black; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_games_player_black ON public.games USING btree (player_black);


--
-- Name: idx_games_player_white; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_games_player_white ON public.games USING btree (player_white);


--
-- Name: idx_players_elo; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_players_elo ON public.players USING btree (elo);


--
-- Name: idx_players_name_norm; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX idx_players_name_norm ON public.players USING btree (name_normalised);


--
-- Name: audit_log_game_affected audit_log_game_affected_audit_log_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_log_game_affected
    ADD CONSTRAINT audit_log_game_affected_audit_log_id_fkey FOREIGN KEY (audit_log_id) REFERENCES public.audit_logs(id);


--
-- Name: audit_log_game_affected audit_log_game_affected_game_ikey_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_log_game_affected
    ADD CONSTRAINT audit_log_game_affected_game_ikey_fkey FOREIGN KEY (game_ikey) REFERENCES public.games(ikey);


--
-- Name: audit_log_player_affected audit_log_player_affected_audit_log_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_log_player_affected
    ADD CONSTRAINT audit_log_player_affected_audit_log_id_fkey FOREIGN KEY (audit_log_id) REFERENCES public.audit_logs(id);


--
-- Name: audit_log_player_affected audit_log_player_affected_player_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_log_player_affected
    ADD CONSTRAINT audit_log_player_affected_player_id_fkey FOREIGN KEY (player_id) REFERENCES public.players(id);


--
-- Name: audit_logs audit_logs_done_by_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.audit_logs
    ADD CONSTRAINT audit_logs_done_by_fkey FOREIGN KEY (done_by) REFERENCES public.admin_users(id);


--
-- Name: games games_player_black_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_player_black_fkey FOREIGN KEY (player_black) REFERENCES public.players(id);


--
-- Name: games games_player_white_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_player_white_fkey FOREIGN KEY (player_white) REFERENCES public.players(id);


--
-- Name: games games_submitter_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.games
    ADD CONSTRAINT games_submitter_fkey FOREIGN KEY (submitter) REFERENCES public.players(id);


--
-- PostgreSQL database dump complete
--


