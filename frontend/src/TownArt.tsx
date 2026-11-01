// Small, local illustrations keep the town charming even with no API key or network.
export function CharacterArt({ id, size = 64 }: { id: string; size?: number }) {
  return (
    <svg
      viewBox="0 0 120 140"
      width={size}
      height={(size * 140) / 120}
      aria-hidden="true"
    >
      <ellipse cx="60" cy="132" rx="35" ry="6" fill="#384a37" opacity=".12" />
      {id === "pip" ? (
        <>
          <path
            d="M83 91 Q120 76 108 48 Q104 76 80 75Z"
            fill="#c57339"
            stroke="#76442e"
            strokeWidth="2"
          />
          <path d="M108 48 Q117 63 108 73L100 65Z" fill="#f5e7c9" />
          <path
            d="M42 121L38 132H52L55 117M68 117L69 132H84L78 120"
            fill="#473f32"
          />
          <path
            d="M40 73Q30 93 36 120Q60 130 85 119L80 76Z"
            fill="#804a42"
            stroke="#593e33"
            strokeWidth="2"
          />
          <path d="M49 75L60 105L70 75" fill="#f0deba" />
          <path d="M56 79L65 80L61 99Z" fill="#c49949" />
          <path
            d="M34 88L22 106L30 111L45 96M82 88L95 75L100 81L90 103"
            fill="#bf753c"
            stroke="#76442e"
            strokeWidth="2"
          />
          <path
            d="M25 17L48 30Q60 25 75 30L96 17L91 54Q90 75 61 83Q30 76 28 53Z"
            fill="#d98c46"
            stroke="#754b32"
            strokeWidth="2"
          />
          <path d="M31 25L46 36L35 43M89 25L76 36L86 42" fill="#754d3c" />
          <path
            d="M30 53Q46 49 60 66Q76 49 91 53Q84 76 61 81Q37 76 30 53"
            fill="#f6e8c9"
          />
          <ellipse cx="45" cy="48" rx="3" ry="4" fill="#353c30" />
          <ellipse cx="77" cy="48" rx="3" ry="4" fill="#353c30" />
          <path d="M55 65Q61 61 67 65L61 72Z" fill="#383c30" />
          <path
            d="M51 73Q61 78 71 73"
            fill="none"
            stroke="#76442e"
            strokeWidth="1.7"
          />
          <circle cx="74" cy="103" r="3" fill="#d3b46c" />
          <path d="M69 110H80" stroke="#c08a65" strokeWidth="2" />
        </>
      ) : id === "momo" ? (
        <>
          <ellipse cx="43" cy="128" rx="15" ry="5" fill="#c18a44" />
          <ellipse cx="78" cy="128" rx="15" ry="5" fill="#c18a44" />
          <path
            d="M38 76Q22 85 19 111Q32 116 43 96M84 76Q103 91 102 112Q88 113 79 96"
            fill="#56808c"
            stroke="#385c66"
            strokeWidth="2"
          />
          <path
            d="M31 80Q28 49 38 32Q58 11 81 31Q94 47 91 81L96 111Q85 130 60 128Q32 130 26 111Z"
            fill="#64929b"
            stroke="#385c66"
            strokeWidth="2"
          />
          <path
            d="M39 49Q48 36 60 53Q71 36 82 49L84 85Q60 103 37 85Z"
            fill="#f4ecd4"
          />
          <path
            d="M30 91L40 75L44 98H77L81 75L91 91L92 117Q60 131 29 116Z"
            fill="#bd984d"
          />
          <path d="M42 103H78V119H42Z" fill="#cba95e" />
          <circle cx="41" cy="96" r="3" fill="#f4e6b2" />
          <circle cx="80" cy="96" r="3" fill="#f4e6b2" />
          <ellipse cx="46" cy="59" rx="3" ry="4" fill="#334d51" />
          <ellipse cx="74" cy="59" rx="3" ry="4" fill="#334d51" />
          <path d="M50 69Q60 61 70 69L60 78Z" fill="#d3a34f" />
          <path
            d="M33 30L40 15Q59 8 79 16L88 32"
            fill="#d7bb76"
            stroke="#a78c53"
            strokeWidth="2"
          />
          <ellipse
            cx="60"
            cy="32"
            rx="43"
            ry="8"
            fill="#e0c88d"
            stroke="#a78c53"
            strokeWidth="2"
          />
          <path
            d="M37 24Q60 30 83 24"
            fill="none"
            stroke="#889258"
            strokeWidth="5"
          />
        </>
      ) : (
        <>
          <path
            d="M40 116L28 130H51L54 118M69 118L72 130H96L81 115"
            fill="#7f985b"
            stroke="#4f6748"
            strokeWidth="2"
          />
          <path
            d="M38 76Q24 86 21 110L30 113L46 94M83 76Q99 90 101 110L92 113L75 94"
            fill="#91a86a"
            stroke="#4f6748"
            strokeWidth="2"
          />
          <path
            d="M39 73L32 114Q57 131 89 114L82 73Z"
            fill="#4b7c70"
            stroke="#365a50"
            strokeWidth="2"
          />
          <path
            d="M46 77V99H75V77M42 106H78V119H42Z"
            fill="#68958a"
            stroke="#355e53"
            strokeWidth="2"
          />
          <path
            d="M28 48Q20 17 42 17Q59 15 61 33Q63 15 82 17Q101 19 93 47Q107 71 83 82Q61 94 35 82Q12 71 28 48Z"
            fill="#91a86a"
            stroke="#4f6748"
            strokeWidth="2"
          />
          <ellipse cx="60" cy="68" rx="29" ry="15" fill="#b9c78a" />
          <path
            d="M39 67Q60 78 82 65"
            fill="none"
            stroke="#4b6243"
            strokeWidth="2"
            strokeLinecap="round"
          />
          <circle cx="41" cy="38" r="12" fill="#e1e2b3" />
          <circle cx="80" cy="38" r="12" fill="#e1e2b3" />
          <ellipse cx="43" cy="39" rx="4" ry="6" fill="#374737" />
          <ellipse cx="78" cy="39" rx="4" ry="6" fill="#374737" />
          <path
            d="M28 24Q59 13 93 25"
            fill="none"
            stroke="#735d3a"
            strokeWidth="6"
          />
          <circle
            cx="44"
            cy="22"
            r="10"
            fill="#b6c7b1"
            stroke="#a98946"
            strokeWidth="4"
          />
          <circle
            cx="77"
            cy="22"
            r="10"
            fill="#b6c7b1"
            stroke="#a98946"
            strokeWidth="4"
          />
          <path d="M54 22H67" stroke="#a98946" strokeWidth="4" />
          <path d="M72 104L78 92" stroke="#dabd76" strokeWidth="4" />
        </>
      )}
    </svg>
  );
}

function Cottage({
  x,
  y,
  scale = 1,
  roof = "#a56e50",
}: {
  x: number;
  y: number;
  scale?: number;
  roof?: string;
}) {
  return (
    <g transform={`translate(${x} ${y}) scale(${scale})`}>
      <path
        d="M0 45L40 12L80 45V108H0Z"
        fill="#eddfb7"
        stroke="#8f8c68"
        strokeWidth="2"
      />
      <path
        d="M-9 46L39 5L90 46L79 51L40 22L0 53Z"
        fill={roof}
        stroke="#775c45"
        strokeWidth="2"
      />
      <path d="M61 21V2H72V31" fill={roof} stroke="#775c45" strokeWidth="2" />
      <path
        d="M32 108V78Q40 68 48 78V108"
        fill="#738770"
        stroke="#576953"
        strokeWidth="2"
      />
      <g fill="#9aaf9a" stroke="#757e5e" strokeWidth="2">
        <rect x="10" y="57" width="15" height="18" rx="2" />
        <rect x="56" y="57" width="15" height="18" rx="2" />
      </g>
      <path
        d="M17 57V76M10 66H26M63 57V76M56 66H72"
        stroke="#e8dcb7"
        strokeWidth="2"
      />
      <path
        d="M5 105H26M54 105H77"
        stroke="#acaa80"
        strokeWidth="5"
        strokeLinecap="round"
      />
    </g>
  );
}

function Tree({ x, y, size = 1 }: { x: number; y: number; size?: number }) {
  return (
    <g transform={`translate(${x} ${y}) scale(${size})`}>
      <path d="M0 0L3 87" stroke="#7f7954" strokeWidth="5" />
      <path d="M2 54L-17 37M2 65L20 47" stroke="#7f7954" strokeWidth="3" />
      <path
        d="M-5-26Q-39-11-35 25Q-46 60-12 65Q1 85 26 60Q49 51 32 23Q38-8 14-17Q6-33-5-26"
        fill="#82966c"
      />
      <path
        d="M0-13L2 67M1 19L-20 5M2 42L22 24"
        fill="none"
        stroke="#b0ba8a"
        strokeWidth="2"
        opacity=".7"
      />
    </g>
  );
}

export function TownArt({ day = 0 }: { day?: number }) {
  return (
    <svg
      className="town-art"
      viewBox="0 0 900 350"
      role="img"
      aria-label="Pip the orange fox, Momo the blue penguin and Zig the green frog in Pebbleton's storybook town square"
    >
      <rect width="900" height="350" fill="#e6e4c6" />
      <circle cx="728" cy="71" r="34" fill="#e6c078" opacity=".85" />
      <g fill="#f3f0d9" opacity=".9">
        <path d="M99 67Q106 49 122 57Q133 33 154 53Q177 48 184 68Z" />
        <path d="M474 51Q490 24 510 42Q529 30 538 48Q553 43 561 57H473Z" />
      </g>
      <path
        d="M0 153Q130 56 284 154Q440 50 576 146Q733 64 900 151V350H0Z"
        fill="#bac4a0"
      />
      <path d="M0 198Q163 130 342 186Q572 103 900 196V350H0Z" fill="#a6b58c" />
      <path d="M0 258Q154 207 318 227Q563 190 900 246V350H0Z" fill="#c9cba1" />
      <path
        d="M492 211Q396 252 413 282Q431 308 340 350H651Q529 295 547 273Q568 246 548 211Z"
        fill="#e2d7b1"
      />
      <g opacity=".8">
        <Cottage x={131} y={119} scale={0.79} />
        <Cottage x={225} y={144} scale={0.68} roof="#85947a" />
        <Cottage x={624} y={133} scale={0.78} roof="#b78460" />
        <Cottage x={704} y={160} scale={0.62} />
      </g>
      <g transform="translate(377 71)">
        <rect
          x="18"
          y="88"
          width="122"
          height="91"
          fill="#e9dcb9"
          stroke="#8b8b68"
          strokeWidth="2"
        />
        <path
          d="M5 91L77 37L151 91Z"
          fill="#986d54"
          stroke="#6d664d"
          strokeWidth="2"
        />
        <path
          d="M60 53V7H94V54"
          fill="#d9cfaa"
          stroke="#8b8b68"
          strokeWidth="2"
        />
        <path
          d="M52 10L77-12L101 10Z"
          fill="#7d8c70"
          stroke="#626e55"
          strokeWidth="2"
        />
        <circle
          cx="77"
          cy="29"
          r="13"
          fill="#eee7c9"
          stroke="#7f8765"
          strokeWidth="2"
        />
        <path
          d="M77 20V29L84 33"
          stroke="#737a59"
          strokeWidth="2"
          fill="none"
        />
        <path d="M61 178V139Q77 124 94 139V178" fill="#74866c" />
        <g fill="#a2b29a" stroke="#828b68" strokeWidth="2">
          <path d="M32 111H49V134H32Z" />
          <path d="M106 111H123V134H106Z" />
        </g>
        <path d="M9 178H148M4 184H153" stroke="#c6c29a" strokeWidth="6" />
      </g>
      <Tree x={75} y={169} size={1.1} />
      <Tree x={807} y={180} size={1.2} />
      <Tree x={320} y={188} size={0.55} />
      <path
        d="M260 142Q448 205 643 148"
        fill="none"
        stroke="#8f9470"
        strokeWidth="2"
      />
      {[0, 1, 2, 3, 4, 5, 6, 7].map((n) => (
        <path
          key={n}
          d={`M${280 + n * 47} ${149 + Math.sin((n / 7) * Math.PI) * 29}l18 4-10 18Z`}
          fill={n % 2 ? "#a47c56" : "#6f8d78"}
        />
      ))}
      <g transform="translate(276 215)">
        <CharacterArt id="pip" size={100} />
      </g>
      <g transform="translate(536 213)">
        <CharacterArt id="momo" size={102} />
      </g>
      <g transform="translate(659 234)">
        <CharacterArt id="zig" size={83} />
      </g>
      <g transform="translate(410 260)">
        <ellipse cx="34" cy="70" rx="53" ry="8" fill="#6c744b" opacity=".13" />
        <path
          d="M0 8L8 68H62L70 8"
          fill="#a98150"
          stroke="#806842"
          strokeWidth="2"
        />
        <ellipse
          cx="35"
          cy="8"
          rx="39"
          ry="11"
          fill="#d2b980"
          stroke="#806842"
          strokeWidth="2"
        />
        <path d="M18 23L23 61M47 23L43 62" stroke="#bd9c64" strokeWidth="3" />
        <path d="M17-4L40-8L53 8L28 13Z" fill="#f2e8c7" stroke="#b3a373" />
        <path d="M26 0L38-2M29 4L43 2" stroke="#a0a17c" strokeWidth="1.5" />
        <circle
          className="town-coin"
          cx="57"
          cy="-4"
          r="9"
          fill="#d7b35d"
          stroke="#ad8a42"
          strokeWidth="2"
        />
        <path d="M57-9V1" stroke="#f4dfa2" strokeWidth="2" />
      </g>
      <g fill="#e8d6a0">
        {[115, 150, 210, 771, 830].map((x) => (
          <g key={x}>
            <path d={`M${x} 320v-17`} stroke="#8a9b69" strokeWidth="2" />
            <circle cx={x} cy="302" r="4" />
            <circle cx={x + 4} cy="307" r="4" />
            <circle cx={x - 4} cy="307" r="4" />
          </g>
        ))}
      </g>
      <path
        d="M580 83q7-8 14 0q7-8 14 0M610 65q5-6 10 0q5-6 10 0"
        fill="none"
        stroke="#7f8b6e"
        strokeWidth="2"
        strokeLinecap="round"
      />
      <g fill="#9ba27d" opacity=".35">
        {[0, 1, 2, 3, 4, 5].map((n) => (
          <ellipse
            key={n}
            cx={455 + (n % 2) * 59}
            cy={318 + n * 5}
            rx="7"
            ry="2"
          />
        ))}
      </g>
      <text
        x="24"
        y="332"
        fill="#7e896b"
        fontSize="9"
        letterSpacing="2"
        fontFamily="Georgia,serif"
      >
        PEBBLETON ·{" "}
        {day === 0 ? "A NEW BEGINNING" : `DAY ${String(day).padStart(2, "0")}`}
      </text>
    </svg>
  );
}
