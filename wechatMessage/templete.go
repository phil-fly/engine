package wechatMessage

const TplPage =
`{
    "msgtype": "markdown",
    "markdown": {
        "content": "{content}"
    }
}
`

const Content =
`{title}
{tplParam}
`


const ParamTable = `>{name}:<font color=\"{color}\">{value}</font>
`