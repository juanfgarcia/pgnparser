{{/*

	This template provides the basic functionality to just show
	each game of a collection in a different page.

	Every page starts with a nice header showing some
	administrative information about the game including players'
	names, their ELO, the winner, ECO, ... 

*/}}

\documentclass{article}

\usepackage[utf8]{inputenc}
\usepackage[english]{babel}
\usepackage{mathpazo}
\usepackage{skak}

\def\hrulefill{\leavevmode\leaders\hrule height 10pt\hfill\kern\z@}

{{/* ----------------------------- Main Body ----------------------------- */}}

\begin{document}

{{/*
	For all games, just show the header and then the moves
	Finally, show a diagram with the final position of the game
*/}}

{{range .GetGames}} 

{{/* ------------------------------- Header ------------------------------ */}}

\begin{center}
  {\Large {{.GetTagValue ("Event")}} ({{.GetTagValue ("TimeControl")}})}
\end{center}

\hrule
\vspace{0.1cm}
\noindent
\raisebox{-5pt}{\WhiteKnightOnWhite} {{.GetTagValue ("White")}} ({{.GetTagValue ("WhiteElo")}}) \hfill {{.GetTagValue ("Date")}}\\
\raisebox{-5pt}{\BlackKnightOnWhite} {{.GetTagValue ("Black")}} ({{.GetTagValue ("BlackElo")}}) \hfill {{.GetTagValue ("ECO")}}
\hrule

\vspace{0.5cm}

{{/* -------------------------------- Moves ------------------------------ */}}

\newgame
{{.GetLaTeXMovesWithComments}}\hfill \textbf{ {{.GetTagValue ("Result")}}}

{{/* --------------------------- Final position -------------------------- */}}

\begin{center}
  \showboard
\end{center}

\clearpage

{{end}}

\end{document}
